package httpheader

const (
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cross-Origin-Resource-Policy
	CrossOriginResourcePolicy = "Cross-Origin-Resource-Policy"

	CrossOriginResourcePolicy_SameSite    CrossOriginResourcePolicyDirective = "same-site"
	CrossOriginResourcePolicy_SameOrigin  CrossOriginResourcePolicyDirective = "same-origin"
	CrossOriginResourcePolicy_CrossOrigin CrossOriginResourcePolicyDirective = "cross-origin"
)

type CrossOriginResourcePolicyDirective string

func (corp CrossOriginResourcePolicyDirective) String() string { return string(corp) }

// https://resourcepolicy.fyi/
func SetCrossOriginResourcePolicy(h Header, policy CrossOriginResourcePolicyDirective) {
	h.Set(CrossOriginResourcePolicy, policy.String())
}
