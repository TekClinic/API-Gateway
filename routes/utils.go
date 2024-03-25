package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-contrib/location"

	"github.com/TekClinic/API-Gateway/schemas"
	"github.com/gin-gonic/gin"
	sf "github.com/sa-/slicefunk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	SkipParameter  = "skip"
	LimitParameter = "limit"
)

// CreateNamedAPIResourceList creates NamedAPIResourceList for the given request.
func CreateNamedAPIResourceList(ctx *gin.Context,
	skip int32, limit int32, count int32, ids []int32) schemas.NamedAPIResourceList {
	previous, next := GetPaginationLinks(ctx, skip, limit, count)
	return schemas.NamedAPIResourceList{
		Count:    count,
		Next:     next,
		Previous: previous,
		Results: sf.Map(ids, func(id int32) schemas.NamedAPIResource {
			return CreateNamedAPIResource(ctx, resourceName, id)
		}),
	}
}

// CreateNamedAPIResource creates NamedAPIResource for resourceName with given id.
func CreateNamedAPIResource(ctx *gin.Context, resourceName string, id int32) schemas.NamedAPIResource {
	requestURL := retrieveRequestURL(ctx)
	requestURL.RawQuery = ""
	requestURL.Path = fmt.Sprintf("/%s/%d", resourceName, id)
	return schemas.NamedAPIResource{
		Name: resourceName,
		URL:  requestURL.String(),
	}
}

// GetPaginationLinks creates previous and next links for pagination.
func GetPaginationLinks(ctx *gin.Context, skip int32, limit int32, count int32) (*string, *string) {
	var previous, next *string
	if skip > 0 {
		previousString := replacePaginationParameters(retrieveRequestURL(ctx), max(0, skip-limit), limit).String()
		previous = &previousString
	}
	if skip+limit < count {
		nextString := replacePaginationParameters(retrieveRequestURL(ctx), skip+limit, limit).String()
		next = &nextString
	}
	return previous, next
}

// retrieveRequestURL return URL that contains all available data about request URL.
func retrieveRequestURL(ctx *gin.Context) *url.URL {
	cloned := *ctx.Request.URL
	locationURL := location.Get(ctx)
	if locationURL != nil {
		cloned.Host = locationURL.Host
		cloned.Scheme = locationURL.Scheme
	}
	return &cloned
}

// replacePaginationParameters replaces skip and limit query parameters of the url with new values.
func replacePaginationParameters(url *url.URL, skip int32, limit int32) *url.URL {
	values := url.Query()

	values.Set(SkipParameter, strconv.Itoa(int(skip)))
	values.Set(LimitParameter, strconv.Itoa(int(limit)))

	url.RawQuery = values.Encode()
	return url
}

// UnImplemented handler for unimplemented endpoints.
func UnImplemented() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusNotImplemented, schemas.ErrorResponse{
			Message: "endpoint is not yet implemented",
		})
	}
}

// HandleGRPCError ends connection with a relevant status code and message.
func HandleGRPCError(err error, ctx *gin.Context) {
	switch status.Code(err) {
	case codes.Unauthenticated:
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, schemas.ErrorResponse{
			Message: "invalid authentication token",
		})
	case codes.PermissionDenied:
		ctx.AbortWithStatusJSON(http.StatusForbidden, schemas.ErrorResponse{
			Message: "you are not allowed to do this",
		})
	case codes.NotFound:
		ctx.AbortWithStatusJSON(http.StatusNotFound, schemas.ErrorResponse{
			Message: "request object is not found",
		})
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: fmt.Sprintf("unknown error occurred: %s", err.Error()),
		})
	}
}
