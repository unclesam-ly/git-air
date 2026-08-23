package updater

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// CurrentVersion 当前编译版本号
const CurrentVersion = "v1.0.7"

const (
	githubAPIURL  = "https://api.github.com/repos/unclesam-ly/git-air/releases/latest"
	checkInterval = 24 * time.Hour
	httpTimeout   = 1500 * time.Millisecond
)

type UpdateState struct {
	LastChecked int64  `json:"last_checked"`
	LatestVer   string `json:"latest_ver"`
	ReleaseURL  string `json:"release_url"`
}

type UpdateInfo struct {
	CurrentVersion string
	LatestVersion  string
	ReleaseURL     string
}

func getCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".git-air")
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "update_check.json"), nil
}

func readCache() (*UpdateState, error) {
	path, err := getCachePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var state UpdateState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func writeCache(state *UpdateState) {
	path, err := getCachePath()
	if err != nil {
		return
	}
	data, err := json.Marshal(state)
	if err == nil {
		_ = os.WriteFile(path, data, 0600)
	}
}

// CheckUpdate 检查是否有新版本（带 24 小时本地缓存与极短超时，静默处理所有异常）
func CheckUpdate(force bool) *UpdateInfo {
	now := time.Now().Unix()

	// 1. 如果不是强制检查，优先读取 24 小时内的本地缓存
	if !force {
		if state, err := readCache(); err == nil && state != nil {
			if now-state.LastChecked < int64(checkInterval.Seconds()) {
				if IsNewerVersion(state.LatestVer, CurrentVersion) {
					return &UpdateInfo{
						CurrentVersion: CurrentVersion,
						LatestVersion:  state.LatestVer,
						ReleaseURL:     state.ReleaseURL,
					}
				}
				return nil
			}
		}
	}

	// 2. 发起轻量级 GitHub API 请求
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", githubAPIURL, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "git-air-cli")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var release struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil || release.TagName == "" {
		return nil
	}

	// 3. 写入本地缓存
	writeCache(&UpdateState{
		LastChecked: now,
		LatestVer:   release.TagName,
		ReleaseURL:  release.HTMLURL,
	})

	// 4. 比对版本
	if IsNewerVersion(release.TagName, CurrentVersion) {
		return &UpdateInfo{
			CurrentVersion: CurrentVersion,
			LatestVersion:  release.TagName,
			ReleaseURL:     release.HTMLURL,
		}
	}

	return nil
}

// CheckAsync 后台并发静默检查，返回 channel 供主流程结束前快速获取
func CheckAsync() <-chan *UpdateInfo {
	ch := make(chan *UpdateInfo, 1)
	go func() {
		info := CheckUpdate(false)
		ch <- info
	}()
	return ch
}

// IsNewerVersion 语义化版本号比对 (如 "v1.0.5" > "v1.0.4")
func IsNewerVersion(latest, current string) bool {
	latestClean := strings.TrimPrefix(strings.TrimSpace(latest), "v")
	currentClean := strings.TrimPrefix(strings.TrimSpace(current), "v")

	if latestClean == "" || currentClean == "" || latestClean == currentClean {
		return false
	}

	latestParts := strings.Split(latestClean, ".")
	currentParts := strings.Split(currentClean, ".")

	for i := 0; i < len(latestParts) && i < len(currentParts); i++ {
		lNum, err1 := strconv.Atoi(latestParts[i])
		cNum, err2 := strconv.Atoi(currentParts[i])
		if err1 != nil || err2 != nil {
			return latestClean > currentClean
		}
		if lNum > cNum {
			return true
		} else if lNum < cNum {
			return false
		}
	}

	return len(latestParts) > len(currentParts)
}
