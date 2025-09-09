package gapi

import (
	"context"
	"fmt"
	"simplebank/token"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context,accessibleRoles []string) (*token.Payload, error) {
	md,ok:=metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) <2 {
		return nil ,fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType!= authorizationBearer {
        return nil, fmt.Errorf("unsupported authorization type %s", authType)
    }

	accessToken := fields[1]
	payload,err :=server.tokenMaker.VerifyToken(accessToken)
	if err!= nil {
        return nil, fmt.Errorf("invalid access token %s", err)
    }

	if !hasPermission(payload.Role,accessibleRoles){
		return nil , fmt.Errorf("permission denied")
	}

	return payload,nil
}

func hasPermission(userRole string, accessibleRoles []string)bool{
	for _, role := range accessibleRoles {
        if role == userRole {
            return true
        }
    }
    return false  // User does not have any of the required roles.  Return false.  Note that this is a simple example and may need to be adjusted depending on your specific needs.  For example, you might want to check if the user's role is in the intersection of the accessible roles.  However, this example assumes that the roles are stored as strings.  In a real-world application, you would need to manage the roles in a database or a similar way.  For example, you could use a slice of strings or a map of strings.  The keys of the map would be the roles, and the values would be booleans indicating whether the user has that role.  You would then check the user's role against the keys in the map.  If the user's role is present in the map, then they have that role.
}