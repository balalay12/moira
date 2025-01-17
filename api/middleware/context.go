package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/moira-alert/moira"
	"github.com/moira-alert/moira/api"
	metricSource "github.com/moira-alert/moira/metric_source"
)

// DatabaseContext sets to requests context configured database
func DatabaseContext(database moira.Database) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := context.WithValue(request.Context(), databaseKey, database)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

// SearchIndexContext sets to requests context configured moira.index.searchIndex
func SearchIndexContext(searcher moira.Searcher) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := context.WithValue(request.Context(), searcherKey, searcher)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

// UserContext get x-webauth-user header and sets it in request context, if header is empty sets empty string
func UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userLogin := request.Header.Get("x-webauth-user")
		ctx := context.WithValue(request.Context(), loginKey, userLogin)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// TriggerContext gets triggerId from parsed URI corresponding to trigger routes and set it to request context
func TriggerContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		triggerID := chi.URLParam(request, "triggerId")
		if triggerID == "" {
			render.Render(writer, request, api.ErrorInvalidRequest(fmt.Errorf("triggerID must be set"))) //nolint
			return
		}
		ctx := context.WithValue(request.Context(), triggerIDKey, triggerID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// ContactContext gets contactID from parsed URI corresponding to trigger routes and set it to request context
func ContactContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		contactID := chi.URLParam(request, "contactId")
		if contactID == "" {
			render.Render(writer, request, api.ErrorInvalidRequest(fmt.Errorf("contactID must be set"))) //nolint
			return
		}
		ctx := context.WithValue(request.Context(), contactIDKey, contactID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// TagContext gets tagName from parsed URI corresponding to tag routes and set it to request context
func TagContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tag := chi.URLParam(request, "tag")
		if tag == "" {
			render.Render(writer, request, api.ErrorInvalidRequest(fmt.Errorf("tag must be set"))) //nolint
			return
		}
		ctx := context.WithValue(request.Context(), tagKey, tag)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// SubscriptionContext gets subscriptionId from parsed URI corresponding to subscription routes and set it to request context
func SubscriptionContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		subscriptionID := chi.URLParam(request, "subscriptionId")
		if subscriptionID == "" {
			render.Render(writer, request, api.ErrorInvalidRequest(fmt.Errorf("subscriptionId must be set"))) //nolint
			return
		}
		ctx := context.WithValue(request.Context(), subscriptionIDKey, subscriptionID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// MetricSourceProvider adds metrics source provider to context
func MetricSourceProvider(sourceProvider *metricSource.SourceProvider) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := context.WithValue(request.Context(), metricSourceProvider, sourceProvider)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

// Paginate gets page and size values from URI query and set it to request context. If query has not values sets given values
func Paginate(defaultPage, defaultSize int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			urlValues, err := url.ParseQuery(request.URL.RawQuery)
			if err != nil {
				render.Render(writer, request, api.ErrorInvalidRequest(err)) //nolint
				return
			}

			page, err := strconv.ParseInt(urlValues.Get("p"), 10, 64)
			if err != nil {
				page = defaultPage
			}

			size, err := strconv.ParseInt(urlValues.Get("size"), 10, 64)
			if err != nil {
				size = defaultSize
			}

			ctxPage := context.WithValue(request.Context(), pageKey, page)
			ctxSize := context.WithValue(ctxPage, sizeKey, size)
			next.ServeHTTP(writer, request.WithContext(ctxSize))
		})
	}
}

// Pager is a function that takes pager id from query
func Pager(defaultCreatePager bool, defaultPagerID string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			urlValues, err := url.ParseQuery(request.URL.RawQuery)
			if err != nil {
				render.Render(writer, request, api.ErrorInvalidRequest(err)) //nolint
				return
			}

			pagerID := urlValues.Get("pagerID")
			if pagerID == "" {
				pagerID = defaultPagerID
			}

			createPager, err := strconv.ParseBool(urlValues.Get("createPager"))
			if err != nil {
				createPager = defaultCreatePager
			}

			ctxPager := context.WithValue(request.Context(), pagerIDKey, pagerID)
			ctxSize := context.WithValue(ctxPager, createPagerKey, createPager)
			next.ServeHTTP(writer, request.WithContext(ctxSize))
		})
	}
}

// Populate gets bool value populate from URI query and set it to request context. If query has not values sets given values
func Populate(defaultPopulated bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			urlValues, err := url.ParseQuery(request.URL.RawQuery)
			if err != nil {
				render.Render(writer, request, api.ErrorInvalidRequest(err)) //nolint
				return
			}

			populate, err := strconv.ParseBool(urlValues.Get("populated"))
			if err != nil {
				populate = defaultPopulated
			}

			ctxTemplate := context.WithValue(request.Context(), populateKey, populate)
			next.ServeHTTP(writer, request.WithContext(ctxTemplate))
		})
	}
}

// Triggers gets string value target from URI query and set it to request context. If query has not values sets given values
func Triggers(localMetricTTL, remoteMetricTTL, prometheusMetricTTL time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()

			ctx = context.WithValue(ctx, localMetricTTLKey, localMetricTTL)
			ctx = context.WithValue(ctx, remoteMetricTTLKey, remoteMetricTTL)
			ctx = context.WithValue(ctx, prometheusMetricTTLKey, prometheusMetricTTL)

			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

// DateRange gets from and to values from URI query and set it to request context. If query has not values sets given values
func DateRange(defaultFrom, defaultTo string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			urlValues, err := url.ParseQuery(request.URL.RawQuery)
			if err != nil {
				render.Render(writer, request, api.ErrorInvalidRequest(err)) //nolint
				return
			}

			from := urlValues.Get("from")
			if from == "" {
				from = defaultFrom
			}

			to := urlValues.Get("to")
			if to == "" {
				to = defaultTo
			}

			ctxPage := context.WithValue(request.Context(), fromKey, from)
			ctxSize := context.WithValue(ctxPage, toKey, to)
			next.ServeHTTP(writer, request.WithContext(ctxSize))
		})
	}
}

// TargetName is a function that gets target name value from query string and places it in context. If query does not have value sets given value.
func TargetName(defaultTargetName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			urlValues, err := url.ParseQuery(request.URL.RawQuery)
			if err != nil {
				render.Render(writer, request, api.ErrorInvalidRequest(err)) //nolint
				return
			}

			targetName := urlValues.Get("target")
			if targetName == "" {
				targetName = defaultTargetName
			}

			ctx := context.WithValue(request.Context(), targetNameKey, targetName)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

// TeamContext gets teamId from parsed URI corresponding to team routes and set it to request context
func TeamContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		teamID := chi.URLParam(request, "teamId")
		if teamID == "" {
			render.Render(writer, request, api.ErrorInvalidRequest(fmt.Errorf("teamId must be set"))) //nolint:errcheck
			return
		}
		ctx := context.WithValue(request.Context(), teamIDKey, teamID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// TeamUserIDContext gets userId from parsed URI corresponding to team routes and set it to request context
func TeamUserIDContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userID := chi.URLParam(request, "teamUserId")
		if userID == "" {
			render.Render(writer, request, api.ErrorInvalidRequest(fmt.Errorf("userId must be set"))) //nolint:errcheck
			return
		}
		ctx := context.WithValue(request.Context(), teamUserIDKey, userID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
