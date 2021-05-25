module github.com/arglucas/secret-app

go 1.16

require (
	"github.com/arglucas/secret-app/handlers" v0.0.0
)

replace (
	github.com/arglucas/secret-app/handlers => ./handlers
)