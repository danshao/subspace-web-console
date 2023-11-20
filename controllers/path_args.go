package controllers

const (
	ARG_USER_ID          = ":user_id"
	ARG_PROFILE_ID       = ":profile_id"
	ARG_PROFILE_PLATFORM = ":platform"
	TYPE_NUMBER          = "([0-9]+)"
	TYPE_STRING          = "([a-zA-Z]+)"
	TYPE_PLATFORM        = "(apple|windows)"
)

var (
	PATH_USER = ARG_USER_ID + TYPE_NUMBER
	PATH_PROFILE = ARG_PROFILE_ID + TYPE_NUMBER
	PATH_PLATFORM = ARG_PROFILE_PLATFORM + TYPE_PLATFORM
)