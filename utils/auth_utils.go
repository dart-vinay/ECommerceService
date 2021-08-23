package utils

func VerifyRequestCredentials(username string, authToken string) error {
	if username=="" || authToken==""{
		return ErrBadRequest
	}
	return nil
}
