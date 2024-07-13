package cookies

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type TypedCookie[T any] struct {
	http.Cookie
	TypedValue T
}

func Marshal[T any](c TypedCookie[T]) (http.Cookie, error) {
	rv := reflect.ValueOf(c.TypedValue)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	v := []string{}
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Type().Field(i)
		fv := rv.FieldByName(f.Name)

		if fv.Kind() == reflect.Pointer {
			fv = fv.Elem()
		}

		if !fv.IsValid() {
			return http.Cookie{}, errors.New("No such field in object: " + f.Name)
		}

		t := f.Tag.Get("cookie")
		if t == "" {
			t = f.Name
		}
		tv := strings.Split(t, ",")

		var value string

		switch fv.Kind() {
		case reflect.Bool:
			value = strconv.FormatBool(fv.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = strconv.FormatInt(fv.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = strconv.FormatUint(fv.Uint(), 10)
		case reflect.Float32:
			value = strconv.FormatFloat(fv.Float(), 'g', -1, 32)
		case reflect.Float64:
			value = strconv.FormatFloat(fv.Float(), 'g', -1, 64)
		case reflect.Complex64:
			value = strconv.FormatComplex(fv.Complex(), 'g', -1, 64)
		case reflect.Complex128:
			value = strconv.FormatComplex(fv.Complex(), 'g', -1, 128)
		case reflect.String:
			value = fv.String()
		default:
			return http.Cookie{},
				errors.New("Unable to convert \"" + value + "\" to field \"" + f.Name + "\" kind \"" + fv.Kind().String() + "\"")
		}

		v = append(v, tv[0]+":"+value)
	}

	return http.Cookie{
		Name:       c.Name,
		Value:      strings.Join(v, "|"),
		Path:       c.Path,
		Domain:     c.Domain,
		Expires:    c.Expires,
		RawExpires: c.RawExpires,
		MaxAge:     c.MaxAge,
		Secure:     c.Secure,
		HttpOnly:   c.HttpOnly,
		SameSite:   c.SameSite,
	}, nil
}

func Unmarshal[T any](data http.Cookie, v *TypedCookie[T]) error {
	if reflect.ValueOf(v).Kind() != reflect.Pointer {
		return errors.New("`v` is not a pointer: " + reflect.ValueOf(v).Kind().String())
	}

	m := make(map[string]string)
	for _, pair := range strings.Split(data.Value, "|") {
		pairV := strings.Split(pair, ":")
		if len(pairV) == 0 {
			return errors.New("Error trying to decode cookie value:\n" + data.Value + "\n\nMissing \":\" pair in first slice")
		}

		key := pairV[0]

		var value string
		if len(pairV) == 1 {
			value = ""
		} else {
			value = strings.Join(pairV[1:], ":")
		}

		m[key] = value
	}

	tt := reflect.TypeOf(&v.TypedValue)
	tv := reflect.ValueOf(&v.TypedValue)
	if tt == nil {
		return errors.New("TypedCookie.TypedValue is not a valid struct{} type")
	}

	if tt.Kind() == reflect.Pointer {
		tt = tt.Elem()
	}
	if tv.Kind() == reflect.Pointer {
		tv = tv.Elem()
	}

	for i := 0; i < tt.NumField(); i++ {
		f := tt.Field(i)
		fv := tv.FieldByName(f.Name)
		if fv.Kind() == reflect.Pointer {
			fv = fv.Elem()
		}

		if !fv.IsValid() {
			return errors.New("No such field in object: " + f.Name)
		}
		if !fv.CanSet() {
			return errors.New("Cannot set value of such field in object: " + f.Name)
		}

		t := f.Tag.Get("cookie")
		if t == "" {
			t = f.Name
		}
		tk := strings.Split(t, ",")[0]

		value := m[tk]

		switch fv.Kind() {
		case reflect.String:
			fv.Set(reflect.ValueOf(value))
		default:
			return errors.New("Unable to convert \"" + value + "\" to field \"" + f.Name + "\" kind \"" + fv.Kind().String() + "\"")
		}
	}
	return nil
}
