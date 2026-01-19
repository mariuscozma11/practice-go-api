# Prequesites

I want to build an app that has the follwing:
- crud api for posts, users with database
- jwt refresh authentication
- RBAC 
- rate limitting
- cors


# Endpoints
## Auth:
|HTTP METHOD|Endpoint           |RBAC |
|--------   |-------------------|-----|
|POST       |/auth/login        |None |

## Posts:
|HTTP METHOD|Endpoint           |RBAC       |
|--------   |-------------------|-----------|
|POST       |/posts             |Admin only |
|PATCH      |/posts/[id]        |Admin only |
|GET        |/posts             |None       |
|GET        |/posts/[id]        |None       |

# Structure convention