package testingutils

func FakeEncryption(data []byte) []byte {
	out := []byte("__fake_encrypted(")
	out = append(out, data...)
	out = append(out, []byte(")")...)
	return out
}

func FakeEcdsaSign(root []byte, address []byte) []byte {
	out := []byte("__fake_ecdsa_sign(root=")
	out = append(out, root...)
	out = append(out, []byte(",address=")...)
	out = append(out, address...)
	out = append(out, []byte(")")...)
	return out
}
