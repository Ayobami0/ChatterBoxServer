package constant

const (
  //WEBSOCKET
  WEBSOCKET_ENDPOINT = "/ws"

	// AUTH
	LOGIN_ENDPOINT  = "/login"
	SIGNUP_ENDPOINT = "/signup"
	LOGOUT_ENDPOINT = "/logout"

	// USER
  USER_GET_ENDPOINT = "/users"
	USER_PROFILE_ENDPOINT = "/users/me"
  USER_CONVERSATONS_ENDPOINT = "/users/me/conversations"
  USER_REQUEST_ENDPOINT = "/users/me/requests"
  USER_REQUEST_ACCEPT_ENDPOINT = "/users/me/requests/:id/accept"
  USER_REQUEST_REJECT_ENDPOINT = "/users/me/requests/:id/reject"
  USER_REQUEST_INVITE_ENDPOINT = "/users/:id/invite"

	// CONVERSATION
	CONVERSATION_CREATE_ENDPOINT  = "/conversations/new"
	CONVERSATION_GET_ENDPOINT     = "/conversations"
	CONVERSATION_DELETE_ENDPOINT  = "/conversations/:id"
	CONVERSATION_JOIN_ENDPOINT    = "/conversations/:id/join"
  CONVERSATION_REQUEST_ENDPOINT = "/conversations/:id/requests"
  CONVERSATION_REQUEST_REJECT = "/conversations/:c_id/requests/:r_id/reject"
  CONVERSATION_REQUEST_ACCEPT = "/conversations/:c_id/requests/:r_id/accept"
)
