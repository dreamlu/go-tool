package lib

import (
	"testing"
)

func TestAesEn(t *testing.T) {

	var as = NewAes()
	t.Log("[aesEn]:", as.AesEn("admin"))
	t.Log("[aesDe]:", as.AesDe("sPa0sTmDf6gasS9tHvIqKw=="))
	t.Log(as.IsAes("13242trergf"))
	t.Log(as.IsAes("sPa0sTmDf6gasS9tHvIqKw=="))
}

func TestRemoveDuplicate(t *testing.T) {
	ss := []string{"a", "b", "c", "a", "b"}
	t.Log(RemoveDuplicate(ss))
	t.Log(Remove(ss, "b"))
}
