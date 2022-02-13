package handlers

import (
	"context"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	elasticCtxKey
)

func CtxLog(entry *logrus.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logrus.Entry {
	return r.Context().Value(logCtxKey).(*logrus.Entry)
}

func CtxElasticCli(q *elastic.Client) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, elasticCtxKey, q)
	}
}

func ElasticCli(r *http.Request) *elastic.Client {
	return r.Context().Value(elasticCtxKey).(*elastic.Client)
}


