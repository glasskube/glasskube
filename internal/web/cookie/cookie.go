package cookie

import (
	"net/http"
	"strconv"
)

const AdvancedOptionsKey = "advancedOptions"

func GetAdvancedOptionsFromCookie(r *http.Request) (bool, error) {
	advancedOptionsVal := false
	if c, err := r.Cookie(AdvancedOptionsKey); err == nil {
		if b, err := strconv.ParseBool(c.Value); err != nil {
			return false, err
		} else {
			advancedOptionsVal = b
		}
	}
	return advancedOptionsVal, nil
}

func SetAdvancedOptionsCookie(w http.ResponseWriter, advancedOptionsEnabled bool) {
	cookie := http.Cookie{
		Name:     AdvancedOptionsKey,
		Value:    strconv.FormatBool(advancedOptionsEnabled),
		MaxAge:   60 * 60 * 24 * 365,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
}
