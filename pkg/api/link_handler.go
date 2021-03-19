package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/fransoaardi/url-shortener/internal/redis"
	"github.com/fransoaardi/url-shortener/internal/utils"
)

type linkHandler struct {
	rc *redis.Client
}

func NewLinkHandler() *linkHandler {
	return &linkHandler{
		rc: redis.NewClient(),
	}
}

func (l *linkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	if id, ok := r.Context().Value("linkID").(string); ok {
		url, err := l.rc.HGet(r.Context(), id, "url")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("error occurred"))
			return
		}
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("no result"))
	}
}

type GenerateRequest struct {
	URL string `json:"url"`
}

type GenerateResponse struct {
	LinkID string `json:"linkId"`
	URL    string `json:"url"`
}

const secondsOneMonth = 60 * 60 * 24 * 30

func (l *linkHandler) Generate(w http.ResponseWriter, r *http.Request) {
	read, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("잘못된 요청 포맷"))
		return
	}

	var gReq GenerateRequest
	if err := json.Unmarshal(read, &gReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("잘못된 요청 포맷"))
		return
	}

	initBase := 0
	for {
		hashed, base, err := utils.GenerateID(gReq.URL, initBase)
		if err != nil {
			switch {
			case errors.Is(err, utils.ErrorHashCollision):
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("hash 실패"))
				return
			case errors.Is(err, utils.ErrorMalformedURL):
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("잘못된 url 입니다"))
				return
			case errors.Is(err, utils.ErrorTryAgain):
				initBase = base + 1 // 다음 base 로 재시도
				continue
			default:
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("예상치못한 error 발생"))
				return
			}
		}

		if saved, _ := l.rc.HGet(r.Context(), hashed, "url"); saved != "" {
			if saved == gReq.URL { // 이미 저장된 값이 있고, 저장된 url 이 동일할때
				gResp := GenerateResponse{
					LinkID: hashed,
					URL:    gReq.URL,
				}

				_ = l.rc.Expire(r.Context(), hashed, secondsOneMonth)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(gResp)
				return
			} else { // 이미 저장된 값이 있고, 저장된 url 이 요청과는 다를때 hash collision 이므로 base 를 올려서 다시 시도한다
				initBase = base + 1
			}
		}

		if err := l.rc.HSet(r.Context(), hashed, "url", gReq.URL); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("redis error, 저장하지 못했습니다"))
			return
		} else {
			gResp := GenerateResponse{
				LinkID: hashed,
				URL:    gReq.URL,
			}

			_ = l.rc.Expire(r.Context(), hashed, secondsOneMonth)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(gResp)
			return
		}
	}
}
