package process

import (
	"fmt"
	"kitu/gs1sscc"
	"regexp"
)

func (k *Krinica) GenerateRecords(start, count int, prefix string) error {
	k.Sscc = make([]string, 0)
	i := start
	stop := start + count - 1
	for {
		pal, err := k.GenerateSSCC(i, prefix)
		if err != nil {
			return fmt.Errorf("generate sscc error %w", err)
		}
		k.Sscc = append(k.Sscc, pal)
		i++
		if i > stop {
			break
		}
	}
	return nil
}

func (k *Krinica) GenerateSSCC(i int, prefix string) (string, error) {
	if len(prefix) < 7 {
		return "", fmt.Errorf("prefix SSCC must be string min len 7")
	}
	// Must be 1–12 digits
	if matched, _ := regexp.MatchString(`^\d{1,12}$`, prefix); !matched {
		return "", fmt.Errorf("invalid SsccPrefix %q: must be 1–12 digits", prefix)
	}
	// Left-zero pad or truncate to 10 digits
	// switch {
	// case len(prefix) < 12:
	// 	prefix = strings.Repeat("0", 12-len(prefix)) + prefix
	// case len(prefix) > 12:
	// 	prefix = prefix[:12]
	// }
	code := fmt.Sprintf("%012.12s", prefix) + fmt.Sprintf("%05d", i)
	sscc, err := gs1sscc.Sscc(code)
	if err != nil {
		return "", fmt.Errorf("sscc returned error for code %w", err)
	}
	return "00" + sscc, nil
}
