package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Validator type qui contient un map d'erreurs de validation
type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid() retourne un "vrai" si les FieldErrors map
// si la case n'est pas vide.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddFieldError() génère un message d'erreur vers FieldErrors map
// et ensuite dans la page web
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField() ajoute une message dans FieldErrors map seulement
// si la validation n'est pas 'ok'
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// Vérifie si la case n'est pas vide
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// Retourne vrai si le nombre de caractères est plus grand que 'n'
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

var NniRx = regexp.MustCompile(`^[a-zA-Z]{1}[0-9]{5}$`)

var passRx = regexp.MustCompile(`^[a-zA-Z0-9.!@#$%^&*()_+-={}[]"';:.,?></]{8,}$`)

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}
