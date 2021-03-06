package pongo2

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func init() {
	RegisterFilter("escape", filterEscape)
	RegisterFilter("safe", filterSafe)

	RegisterFilter("add", filterAdd)
	RegisterFilter("capfirst", filterCapfirst)
	RegisterFilter("cut", filterCut)
	RegisterFilter("date", filterDate)
	RegisterFilter("default", filterDefault)
	RegisterFilter("default_if_none", filterDefaultIfNone)
	RegisterFilter("divisibleby", filterDivisibleby)
	RegisterFilter("first", filterFirst)
	RegisterFilter("last", filterLast)
	RegisterFilter("length", filterLength)
	RegisterFilter("length_is", filterLengthis)
	RegisterFilter("linebreaksbr", filterLinebreaksbr)
	RegisterFilter("lower", filterLower)
	RegisterFilter("pluralize", filterPluralize)
	RegisterFilter("removetags", filterRemovetags)
	RegisterFilter("upper", filterUpper)
	RegisterFilter("urlencode", filterUrlencode)
	RegisterFilter("striptags", filterStriptags)
	RegisterFilter("time", filterDate) // time uses filterDate (same golang-format)
	RegisterFilter("truncatechars", filterTruncatechars)
	RegisterFilter("yesno", filterYesno)

	RegisterFilter("float", filterFloat)     // pongo-specific
	RegisterFilter("integer", filterInteger) // pongo-specific

	/* Missing filters:

	   addslashes
	   center
	   dictsort
	   dictsortreversed
	   escape
	   escapejs
	   filesizeformat
	   floatformat
	   force_escape
	   get_digit
	   iriencode
	   join
	   linebreaks
	   linenumbers
	   ljust
	   make_list
	   phone2numeric
	   pprint
	   random
	   removetags
	   rjust
	   safeseq
	   slice
	   slugify
	   stringformat
	   timesince
	   timeuntil
	   title
	   truncatechars_html
	   truncatewords
	   truncatewords_html
	   unordered_list
	   urlize
	   urlizetrunc
	   wordcount
	   wordwrap

	   Filters that won't be added:

	   static
	   get_static_prefix
	*/
}

func filterTruncatechars(in *Value, param *Value) (*Value, error) {
	if !in.CanSlice() {
		return nil, errors.New("")
	}
	return in.Slice(0, param.Integer()), nil
}

func filterEscape(in *Value, param *Value) (*Value, error) {
	output := strings.Replace(in.String(), "&", "&amp;", -1)
	output = strings.Replace(output, ">", "&gt;", -1)
	output = strings.Replace(output, "<", "&lt;", -1)
	output = strings.Replace(output, "\"", "&quot;", -1)
	output = strings.Replace(output, "'", "&#39;", -1)
	return AsValue(output), nil
}

func filterSafe(in *Value, param *Value) (*Value, error) {
	return in, nil // nothing to do here, just to keep track of the safe application
}

func filterAdd(in *Value, param *Value) (*Value, error) {
	if in.IsNumber() && param.IsNumber() {
		if in.IsFloat() || param.IsFloat() {
			return AsValue(in.Float() + param.Float()), nil
		} else {
			return AsValue(in.Integer() + param.Integer()), nil
		}
	}
	// If in/param is not a number, we're relying on the
	// Value's String() convertion and just add them both together
	return AsValue(in.String() + param.String()), nil
}

func filterCut(in *Value, param *Value) (*Value, error) {
	return AsValue(strings.Replace(in.String(), param.String(), "", -1)), nil
}

func filterLength(in *Value, param *Value) (*Value, error) {
	return AsValue(in.Len()), nil
}

func filterLengthis(in *Value, param *Value) (*Value, error) {
	return AsValue(in.Len() == param.Integer()), nil
}

func filterDefault(in *Value, param *Value) (*Value, error) {
	if !in.IsTrue() {
		return param, nil
	}
	return in, nil
}

func filterDefaultIfNone(in *Value, param *Value) (*Value, error) {
	if in.IsNil() {
		return param, nil
	}
	return in, nil
}

func filterDivisibleby(in *Value, param *Value) (*Value, error) {
	if param.Integer() == 0 {
		return AsValue(false), nil
	}
	return AsValue(in.Integer()%param.Integer() == 0), nil
}

func filterFirst(in *Value, param *Value) (*Value, error) {
	if in.CanSlice() {
		return in.Slice(0, 1), nil
	}
	return AsValue(""), nil
}

func filterLast(in *Value, param *Value) (*Value, error) {
	if in.CanSlice() {
		return in.Slice(in.Len()-1, in.Len()), nil
	}
	return AsValue(""), nil
}

func filterUpper(in *Value, param *Value) (*Value, error) {
	return AsValue(strings.ToUpper(in.String())), nil
}

func filterLower(in *Value, param *Value) (*Value, error) {
	return AsValue(strings.ToLower(in.String())), nil
}

func filterCapfirst(in *Value, param *Value) (*Value, error) {
	return AsValue(strings.Title(in.String())), nil
}

func filterDate(in *Value, param *Value) (*Value, error) {
	t, is_time := in.Interface().(time.Time)
	if !is_time {
		return nil, errors.New("Filter input argument must be of type 'time.Time'.")
	}
	return AsValue(t.Format(param.String())), nil
}

func filterFloat(in *Value, param *Value) (*Value, error) {
	return AsValue(in.Float()), nil
}

func filterInteger(in *Value, param *Value) (*Value, error) {
	return AsValue(in.Integer()), nil
}

func filterLinebreaksbr(in *Value, param *Value) (*Value, error) {
	return AsValue(strings.Replace(in.String(), "\n", "<br />", -1)), nil
}

func filterUrlencode(in *Value, param *Value) (*Value, error) {
	return AsValue(url.QueryEscape(in.String())), nil
}

var re_striptags = regexp.MustCompile("<[^>]*?>")

func filterStriptags(in *Value, param *Value) (*Value, error) {
	s := in.String()

	// Strip all tags
	s = re_striptags.ReplaceAllString(s, "")

	return AsValue(strings.TrimSpace(s)), nil
}

func filterPluralize(in *Value, param *Value) (*Value, error) {
	if in.IsNumber() {
		// Works only on numbers
		if param.Len() > 0 {
			endings := strings.Split(param.String(), ",")
			if len(endings) > 2 {
				return nil, errors.New("You cannot pass more than 2 arguments to filter 'pluralize'.")
			}
			if len(endings) == 1 {
				// 1 argument
				if in.Integer() != 1 {
					return AsValue(endings[0]), nil
				}
			} else {
				if in.Integer() != 1 {
					// 2 arguments
					return AsValue(endings[1]), nil
				}
				return AsValue(endings[0]), nil
			}
		} else {
			if in.Integer() != 1 {
				// return default 's'
				return AsValue("s"), nil
			}
		}

		return AsValue(""), nil
	} else {
		return nil, errors.New("Filter 'pluralize' does only work on numbers.")
	}
}

func filterRemovetags(in *Value, param *Value) (*Value, error) {
	s := in.String()
	tags := strings.Split(param.String(), ",")

	// Strip only specific tags
	for _, tag := range tags {
		re := regexp.MustCompile(fmt.Sprintf("</?%s/?>", tag))
		s = re.ReplaceAllString(s, "")
	}

	return AsValue(strings.TrimSpace(s)), nil
}

func filterYesno(in *Value, param *Value) (*Value, error) {
	choices := map[int]string{
		0: "yes",
		1: "no",
		2: "maybe",
	}
	param_string := param.String()
	custom_choices := strings.Split(param_string, ",")
	if len(param_string) > 0 {
		if len(custom_choices) > 3 {
			return nil, errors.New(fmt.Sprintf("You cannot pass more than 3 options to the 'yesno'-filter (got: '%s').", param_string))
		}
		if len(custom_choices) < 2 {
			return nil, errors.New(fmt.Sprintf("You must pass either no or at least 2 arguments to the 'yesno'-filter (got: '%s').", param_string))
		}

		// Map to the options now
		choices[0] = custom_choices[0]
		choices[1] = custom_choices[1]
		if len(custom_choices) == 3 {
			choices[2] = custom_choices[2]
		}
	}

	// maybe
	if in.IsNil() {
		return AsValue(choices[2]), nil
	}

	// yes
	if in.IsTrue() {
		return AsValue(choices[0]), nil
	}

	// no
	return AsValue(choices[1]), nil
}
