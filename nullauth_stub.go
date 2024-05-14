//go:build !nullauth

package daisy

func NullAuth(provider AuthProvider) AuthProvider {
	return provider
}
