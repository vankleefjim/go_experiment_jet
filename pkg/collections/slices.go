package collections

func Map[F any,T any](from []F, f func(F)T) []T{
	out := make([]T,len(from))
	for i,v := range from{
		out[i] = f(v)
	}
	return out
}