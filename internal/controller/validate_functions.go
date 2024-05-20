package controller

import (
	"fmt"
	"github.com/crmin/rathole-operator/api/v1alpha1"
	"strings"
)

func ValidateServer(r *v1alpha1.RatholeServer) error {
	bindAddr := r.Spec.BindAddr
	if len(strings.Split(bindAddr, ":")) != 2 {
		return fmt.Errorf("bindAddr must be in the form of ip:port")
	}

	defaultTokenSet := r.Spec.DefaultToken != ""
	defaultTokenFromSet := r.Spec.DefaultTokenFrom.ConfigMapRef.Name != "" || r.Spec.DefaultTokenFrom.SecretRef.Name != ""
	if defaultTokenSet && defaultTokenFromSet {
		return fmt.Errorf("default token and default token from cannot be set at the same time")
	}

	switch r.Spec.Transport.Type {
	case "tls":
		if r.Spec.Transport.TLS == nil {
			return fmt.Errorf("tls transport requires .spec.transport.tls configuration")
		}

		pkcs12Set := r.Spec.Transport.TLS.PKCS12 != ""
		pkcs12FromSet := r.Spec.Transport.TLS.PKCS12From.ConfigMapRef.Name != "" || r.Spec.Transport.TLS.PKCS12From.SecretRef.Name != ""
		if pkcs12Set && pkcs12FromSet {
			return fmt.Errorf(".spec.transport.tls.pkcs12 and .spec.transport.tls.pkcs12From cannot be set at the same time")
		} else if !pkcs12Set && !pkcs12FromSet {
			return fmt.Errorf("tls transport requires .spec.transport.tls.pkcs12 or .spec.transport.tls.pkcs12From")
		}

		pkcs12PasswordSet := r.Spec.Transport.TLS.PKCS12Password != ""
		pkcs12PasswordFromSet := r.Spec.Transport.TLS.PKCS12PasswordFrom.ConfigMapRef.Name != "" || r.Spec.Transport.TLS.PKCS12PasswordFrom.SecretRef.Name != ""
		if pkcs12PasswordSet && pkcs12PasswordFromSet {
			return fmt.Errorf(".spec.transport.tls.pkcs12Password and .spec.transport.tls.pkcs12PasswordFrom cannot be set at the same time")
		} else if !pkcs12PasswordSet && !pkcs12PasswordFromSet {
			return fmt.Errorf("tls transport requires .spec.transport.tls.pkcs12Password or .spec.transport.tls.pkcs12PasswordFrom")
		}
	case "noise":
		if r.Spec.Transport.Noise == nil {
			return fmt.Errorf("noise transport requires .spec.transport.noise configuration")
		}

		// If set private key, need to set public key too and vice versa for auth
		localPrivateKeySet := r.Spec.Transport.Noise.LocalPrivateKey != ""
		localPrivateKeyFromSet := r.Spec.Transport.Noise.LocalPrivateKeyFrom.ConfigMapRef.Name != "" || r.Spec.Transport.Noise.LocalPrivateKeyFrom.SecretRef.Name != ""
		remotePublicKeySet := r.Spec.Transport.Noise.RemotePublicKey != ""
		remotePublicKeyFromSet := r.Spec.Transport.Noise.RemotePublicKeyFrom.ConfigMapRef.Name != "" || r.Spec.Transport.Noise.RemotePublicKeyFrom.SecretRef.Name != ""

		localPrivateSet := localPrivateKeySet || localPrivateKeyFromSet
		remotePublicSet := remotePublicKeySet || remotePublicKeyFromSet

		if localPrivateKeySet && localPrivateKeyFromSet {
			return fmt.Errorf(".spec.transport.noise.localPrivateKey and .spec.transport.noise.localPrivateKeyFrom cannot be set at the same time")
		}
		if remotePublicKeySet && remotePublicKeyFromSet {
			return fmt.Errorf(".spec.transport.noise.remotePublicKey and .spec.transport.noise.remotePublicKeyFrom cannot be set at the same time")
		}
		if localPrivateSet && !remotePublicSet {
			return fmt.Errorf("noise transport requires .spec.transport.noise.remotePublicKey")
		}
		if !localPrivateSet && remotePublicSet {
			return fmt.Errorf("noise transport requires .spec.transport.noise.localPrivateKey")
		}
	}
	return nil
}

func ValidateClient(r *v1alpha1.RatholeClient) error {
	defaultTokenSet := r.Spec.DefaultToken != ""
	defaultTokenFromSet := r.Spec.DefaultTokenFrom.ConfigMapRef.Name != "" || r.Spec.DefaultTokenFrom.SecretRef.Name != ""
	if defaultTokenSet && defaultTokenFromSet {
		return fmt.Errorf("default token and default token from cannot be set at the same time")
	}

	switch r.Spec.Transport.Type {
	case "tls":
		if r.Spec.Transport.TLS == nil {
			return fmt.Errorf("tls transport requires .spec.transport.tls configuration")
		}
		trustedRootFromSet := r.Spec.Transport.TLS.TrustedRootFrom.ConfigMapRef.Name != "" || r.Spec.Transport.TLS.TrustedRootFrom.SecretRef.Name != ""
		trustedRootSet := r.Spec.Transport.TLS.TrustedRoot != ""
		if trustedRootSet && trustedRootFromSet {
			return fmt.Errorf(".spec.transport.tls.trustedRoot and .spec.transport.tls.trustedRootFrom cannot be set at the same time")
		} else if !trustedRootSet && !trustedRootFromSet {
			return fmt.Errorf("tls transport requires .spec.transport.tls.trustedRoot or .spec.transport.tls.trustedRootFrom")
		}
	case "noise":
		if r.Spec.Transport.Noise == nil {
			return fmt.Errorf("noise transport requires .spec.transport.noise configuration")
		}

		// If set private key, need to set public key too and vice versa for auth
		localPrivateKeySet := r.Spec.Transport.Noise.LocalPrivateKey != ""
		localPrivateKeyFromSet := r.Spec.Transport.Noise.LocalPrivateKeyFrom.ConfigMapRef.Name != "" || r.Spec.Transport.Noise.LocalPrivateKeyFrom.SecretRef.Name != ""
		remotePublicKeySet := r.Spec.Transport.Noise.RemotePublicKey != ""
		remotePublicKeyFromSet := r.Spec.Transport.Noise.RemotePublicKeyFrom.ConfigMapRef.Name != "" || r.Spec.Transport.Noise.RemotePublicKeyFrom.SecretRef.Name != ""

		localPrivateSet := localPrivateKeySet || localPrivateKeyFromSet
		remotePublicSet := remotePublicKeySet || remotePublicKeyFromSet

		if localPrivateKeySet && localPrivateKeyFromSet {
			return fmt.Errorf(".spec.transport.noise.localPrivateKey and .spec.transport.noise.localPrivateKeyFrom cannot be set at the same time")
		}
		if remotePublicKeySet && remotePublicKeyFromSet {
			return fmt.Errorf(".spec.transport.noise.remotePublicKey and .spec.transport.noise.remotePublicKeyFrom cannot be set at the same time")
		}
		if localPrivateSet && !remotePublicSet {
			return fmt.Errorf("noise transport requires .spec.transport.noise.remotePublicKey")
		}
		if !localPrivateSet && remotePublicSet {
			return fmt.Errorf("noise transport requires .spec.transport.noise.localPrivateKey")
		}
	}

	return nil
}

func ValidateService(r *v1alpha1.RatholeService, server *v1alpha1.RatholeServer, client *v1alpha1.RatholeClient) error {
	// Note: Need check server and client exist?

	tokenSet := r.Spec.Token != ""
	clientDefaultTokenSet := client.Spec.DefaultToken != ""
	if !tokenSet && !clientDefaultTokenSet {
		return fmt.Errorf(".spec.token must be set if client .spec.defaultToken was not set")
	}

	return nil
}
