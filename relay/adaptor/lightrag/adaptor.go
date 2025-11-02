package lightrag

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/songquanpeng/one-api/relay/meta"
	"github.com/songquanpeng/one-api/relay/relaymode"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/relay/adaptor"
	"github.com/songquanpeng/one-api/relay/model"
)

type Adaptor struct {
	apiKey  string
	baseUrl string
	uri     *url.URL
}

func (a *Adaptor) Init(meta *meta.Meta) {
	a.apiKey = meta.APIKey
	a.baseUrl = meta.BaseURL
	a.uri, _ = url.Parse(meta.BaseURL)
}

func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	// https://github.com/ollama/ollama/blob/main/docs/api.md
	fullRequestURL := a.uri.JoinPath("/query/stream")
	return fullRequestURL.String(), nil
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) error {
	adaptor.SetupCommonRequestHeader(c, req, meta)
	req.Header.Set("X-Openai-Key", a.apiKey)
	return nil
}

func (a *Adaptor) ConvertRequest(c *gin.Context, relayMode int, request *model.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	switch relayMode {
	case relaymode.Embeddings:
		ollamaEmbeddingRequest := ConvertEmbeddingRequest(*request)
		return ollamaEmbeddingRequest, nil
	default:
		return ConvertRequest(*request), nil
	}
}

func (a *Adaptor) ConvertImageRequest(request *model.ImageRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	return request, nil
}

func (a *Adaptor) DoRequest(c *gin.Context, meta *meta.Meta, requestBody io.Reader) (*http.Response, error) {
	return adaptor.DoRequestHelper(a, c, meta, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, meta *meta.Meta) (usage *model.Usage, err *model.ErrorWithStatusCode) {
	if meta.IsStream {
		err, usage = StreamHandler(c, resp)
	} else {
		switch meta.Mode {
		case relaymode.Embeddings:
			err, usage = EmbeddingHandler(c, resp)
		default:
			err, usage = Handler(c, resp)
		}
	}
	return
}

func (a *Adaptor) GetModelList() []string {
	return ModelList
}

func (a *Adaptor) GetChannelName() string {
	return "Lightrag"
}
