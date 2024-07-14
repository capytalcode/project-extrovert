package cookies

import (
	"errors"
	"math/bits"
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

		value, err := encodeString(fv)
		if err != nil {
			return http.Cookie{}, errors.Join(errors.New("Unsupportted type in struct: "+fv.Kind().String()), err)
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

		final, err := decodeString(m[tk], fv.Kind())
		if err != nil {
			return errors.Join(errors.New("Unsupportted type in struct: "+fv.Kind().String()), err)
		}

		if strings.Contains(strings.ToLower(fv.Kind().String()), "complex") {
			fv.SetComplex(final.(complex128))
		} else if strings.Contains(strings.ToLower(fv.Kind().String()), "uint") {
			fv.SetUint(final.(uint64))
		} else if strings.Contains(strings.ToLower(fv.Kind().String()), "int") {
			fv.SetInt(final.(int64))
		} else {
			fv.Set(reflect.ValueOf(final))
		}
	}
	return nil
}

func encodeString(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), nil

	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'g', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'g', -1, 64), nil

	case reflect.Complex64:
		return strconv.FormatComplex(v.Complex(), 'g', -1, 64), nil
	case reflect.Complex128:
		return strconv.FormatComplex(v.Complex(), 'g', -1, 128), nil

	case reflect.String:
		return v.String(), nil

	default:
		return "", errors.ErrUnsupported
	}
}

func decodeString(v string, k reflect.Kind) (any, error) {
	var final any
	var err error

	switch k {
	case reflect.Bool:
		final, err = strconv.ParseBool(v)

	case reflect.Int8:
		final, err = strconv.ParseInt(v, 10, 8)
	case reflect.Int16:
		final, err = strconv.ParseInt(v, 10, 16)
	case reflect.Int32:
		final, err = strconv.ParseInt(v, 10, 32)
	case reflect.Int64:
		final, err = strconv.ParseInt(v, 10, 64)
	case reflect.Int:
		final, err = strconv.ParseInt(v, 10, bits.UintSize)

	case reflect.Uint8:
		final, err = strconv.ParseUint(v, 10, 8)
	case reflect.Uint16:
		final, err = strconv.ParseUint(v, 10, 16)
	case reflect.Uint32:
		final, err = strconv.ParseUint(v, 10, 32)
	case reflect.Uint64:
		final, err = strconv.ParseUint(v, 10, 64)
	case reflect.Uint:
		final, err = strconv.ParseUint(v, 10, bits.UintSize)

	case reflect.Float32:
		final, err = strconv.ParseFloat(v, 32)
	case reflect.Float64:
		final, err = strconv.ParseFloat(v, 64)

	case reflect.Complex64:
		final, err = strconv.ParseComplex(v, 64)
	case reflect.Complex128:
		final, err = strconv.ParseComplex(v, 128)

	case reflect.String:
		final = v

	default:
		return nil, errors.ErrUnsupported
	}

	return final, err
}
