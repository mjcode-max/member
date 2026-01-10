package face

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"member-pre/internal/infrastructure/config"
	"member-pre/pkg/logger"
)

// FaceService 人脸识别服务接口
type FaceService interface {
	// RegisterFace 注册人脸，返回face_id
	RegisterFace(imageData []byte, externalImageId string) (faceId string, error error)
	// VerifyFace 验证人脸（1:1比对），返回相似度
	VerifyFace(faceId string, imageData []byte) (similarity float64, error error)
}

// HuaweiFRSService 华为云FRS服务实现
type HuaweiFRSService struct {
	config *config.HuaweiFRSConfig
	logger logger.Logger
	client *http.Client
}

// NewHuaweiFRSService 创建华为云FRS服务实例
func NewHuaweiFRSService(cfg *config.HuaweiFRSConfig, log logger.Logger) FaceService {
	return &HuaweiFRSService{
		config: cfg,
		logger: log,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RegisterFace 注册人脸
func (s *HuaweiFRSService) RegisterFace(imageData []byte, externalImageId string) (string, error) {
	s.logger.Info("开始注册人脸", logger.NewField("external_image_id", externalImageId))

	// 将图片转换为base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)

	// 构建请求体
	requestBody := map[string]interface{}{
		"image_base64": imageBase64,
	}
	if externalImageId != "" {
		requestBody["external_image_id"] = externalImageId
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		s.logger.Error("序列化请求体失败", logger.NewField("error", err.Error()))
		return "", fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 构建API URL
	url := fmt.Sprintf("%s/v2/%s/face-sets/%s/faces", s.config.Endpoint, s.config.ProjectID, s.config.FaceSetName)

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		s.logger.Error("创建HTTP请求失败", logger.NewField("error", err.Error()))
		return "", fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Sdk-Date", time.Now().UTC().Format("20060102T150405Z"))

	// 添加AK/SK签名
	if err := s.addSignature(req, bodyBytes); err != nil {
		s.logger.Error("添加签名失败", logger.NewField("error", err.Error()))
		return "", fmt.Errorf("添加签名失败: %w", err)
	}

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("发送HTTP请求失败", logger.NewField("error", err.Error()))
		return "", fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("读取响应失败", logger.NewField("error", err.Error()))
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("华为云FRS API返回错误",
			logger.NewField("status_code", resp.StatusCode),
			logger.NewField("response", string(respBody)),
		)
		return "", fmt.Errorf("华为云FRS API返回错误: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var response struct {
		Faces []struct {
			FaceID string `json:"face_id"`
		} `json:"faces"`
		ErrorCode string `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		s.logger.Error("解析响应失败", logger.NewField("error", err.Error()))
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查是否有错误
	if response.ErrorCode != "" {
		s.logger.Error("华为云FRS API返回业务错误",
			logger.NewField("error_code", response.ErrorCode),
			logger.NewField("error_msg", response.ErrorMsg),
		)
		return "", fmt.Errorf("华为云FRS API错误: %s - %s", response.ErrorCode, response.ErrorMsg)
	}

	// 提取face_id
	if len(response.Faces) == 0 {
		s.logger.Error("未检测到人脸")
		return "", fmt.Errorf("未检测到人脸")
	}

	faceID := response.Faces[0].FaceID
	s.logger.Info("人脸注册成功", logger.NewField("face_id", faceID))
	return faceID, nil
}

// VerifyFace 验证人脸（1:1比对）
func (s *HuaweiFRSService) VerifyFace(faceId string, imageData []byte) (float64, error) {
	s.logger.Info("开始验证人脸", logger.NewField("face_id", faceId))

	// 将图片转换为base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)

	// 构建请求体
	requestBody := map[string]interface{}{
		"face_id":      faceId,
		"image_base64": imageBase64,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		s.logger.Error("序列化请求体失败", logger.NewField("error", err.Error()))
		return 0, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 构建API URL
	url := fmt.Sprintf("%s/v2/%s/face-compare", s.config.Endpoint, s.config.ProjectID)

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		s.logger.Error("创建HTTP请求失败", logger.NewField("error", err.Error()))
		return 0, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Sdk-Date", time.Now().UTC().Format("20060102T150405Z"))

	// 添加AK/SK签名
	if err := s.addSignature(req, bodyBytes); err != nil {
		s.logger.Error("添加签名失败", logger.NewField("error", err.Error()))
		return 0, fmt.Errorf("添加签名失败: %w", err)
	}

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("发送HTTP请求失败", logger.NewField("error", err.Error()))
		return 0, fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("读取响应失败", logger.NewField("error", err.Error()))
		return 0, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("华为云FRS API返回错误",
			logger.NewField("status_code", resp.StatusCode),
			logger.NewField("response", string(respBody)),
		)
		return 0, fmt.Errorf("华为云FRS API返回错误: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var response struct {
		Similarity float64 `json:"similarity"`
		ErrorCode  string  `json:"error_code"`
		ErrorMsg   string  `json:"error_msg"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		s.logger.Error("解析响应失败", logger.NewField("error", err.Error()))
		return 0, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查是否有错误
	if response.ErrorCode != "" {
		s.logger.Error("华为云FRS API返回业务错误",
			logger.NewField("error_code", response.ErrorCode),
			logger.NewField("error_msg", response.ErrorMsg),
		)
		return 0, fmt.Errorf("华为云FRS API错误: %s - %s", response.ErrorCode, response.ErrorMsg)
	}

	s.logger.Info("人脸验证成功", logger.NewField("similarity", response.Similarity))
	return response.Similarity, nil
}

// addSignature 添加AK/SK签名到请求头
func (s *HuaweiFRSService) addSignature(req *http.Request, body []byte) error {
	// 获取时间戳
	timestamp := req.Header.Get("X-Sdk-Date")
	if timestamp == "" {
		timestamp = time.Now().UTC().Format("20060102T150405Z")
		req.Header.Set("X-Sdk-Date", timestamp)
	}

	// 构建待签名字符串
	canonicalRequest := s.buildCanonicalRequest(req, body)
	stringToSign := s.buildStringToSign(timestamp, canonicalRequest)

	// 计算签名
	signature := s.calculateSignature(stringToSign)

	// 设置Authorization头
	authHeader := fmt.Sprintf("SDK-HMAC-SHA256 Access=%s, SignedHeaders=host;x-sdk-date, Signature=%s",
		s.config.AccessKeyID, signature)
	req.Header.Set("Authorization", authHeader)

	return nil
}

// buildCanonicalRequest 构建规范请求
func (s *HuaweiFRSService) buildCanonicalRequest(req *http.Request, body []byte) string {
	// HTTP方法
	method := req.Method

	// URI路径
	uri := req.URL.Path
	if req.URL.RawQuery != "" {
		uri += "?" + req.URL.RawQuery
	}

	// 查询字符串（已包含在URI中，这里为空）
	queryString := ""

	// 请求头（只包含host和x-sdk-date）
	headers := fmt.Sprintf("host:%s\nx-sdk-date:%s\n", req.URL.Host, req.Header.Get("X-Sdk-Date"))
	signedHeaders := "host;x-sdk-date"

	// 请求体哈希
	bodyHash := fmt.Sprintf("%x", sha256.Sum256(body))

	// 组合规范请求
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		method, uri, queryString, headers, signedHeaders, bodyHash)
}

// buildStringToSign 构建待签名字符串
func (s *HuaweiFRSService) buildStringToSign(timestamp, canonicalRequest string) string {
	algorithm := "SDK-HMAC-SHA256"
	canonicalRequestHash := fmt.Sprintf("%x", sha256.Sum256([]byte(canonicalRequest)))

	return fmt.Sprintf("%s\n%s\n%s", algorithm, timestamp, canonicalRequestHash)
}

// calculateSignature 计算签名
func (s *HuaweiFRSService) calculateSignature(stringToSign string) string {
	// 从时间戳提取日期
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	date := timestamp[:8]

	// 计算kDate
	kDate := hmacSHA256([]byte(s.config.SecretAccessKey), date)

	// 计算kService
	kService := hmacSHA256(kDate, "face")

	// 计算kSigning
	kSigning := hmacSHA256(kService, "sdk_request")

	// 计算最终签名
	signature := hmacSHA256(kSigning, stringToSign)

	return fmt.Sprintf("%x", signature)
}

// hmacSHA256 计算HMAC-SHA256
func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

