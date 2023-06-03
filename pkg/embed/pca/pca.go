package pca

func EmbeddingsTo3D(embeddings [][1536]float32) [][]float32 {
	// Calculate the covariance matrix of the input embeddings
	n := len(embeddings)
	m := len(embeddings[0])
	means := [1536]float32(make([]float32, m))
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			means[i] += embeddings[j][i]
		}
		means[i] /= float32(n)
	}
	cov := make([][1536]float32, m)
	for i := 0; i < m; i++ {
		cov[i] = [1536]float32(make([]float32, m))
		for j := 0; j < m; j++ {
			for k := 0; k < n; k++ {
				cov[i][j] += (embeddings[k][i] - means[i]) * (embeddings[k][j] - means[j])
			}
			cov[i][j] /= float32(n - 1)
		}
	}

	// Calculate the eigenvectors and eigenvalues of the covariance matrix
	eigenVals := [1536]float32(make([]float32, m))
	tmpVecs := make([][1536]float32, m)
	for i := 0; i < m; i++ {
		tmpVecs[i] = [1536]float32(make([]float32, m))
	}
	for i := range eigenVals {
		eigenVals[i] = cov[i][i]
		tmpVecs[i][i] = 1.0
	}
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			c := cov[i][j] / ((cov[i][i] + cov[j][j]) / 2)
			s := float32(1.0) / float32(2+2*c*c)
			tmpVecs[i][j] = s * float32(2*c)
			tmpVecs[j][i] = s * float32(-2*c)
			tmp := tmpVecs[i][i]
			tmpVecs[i][i] = tmp*s*(cov[i][i]-cov[j][j]) + tmpVecs[i][j]*c*2*cov[i][j]
			tmpVecs[j][j] = tmp*s*(cov[j][j]-cov[i][i]) + tmpVecs[i][j]*c*2*cov[i][j]
			tmpVecs[i][j] = tmpVecs[j][i]*(cov[i][i]-cov[j][j]) + tmpVecs[i][j]*(cov[i][j]+cov[i][j])
		}
	}
	for i := 0; i < m-1; i++ {
		maxIdx := i
		max := eigenVals[i]
		for j := i; j < m; j++ {
			if eigenVals[j] > max {
				maxIdx = j
				max = eigenVals[j]
			}
		}
		if maxIdx != i {
			eigenVals[i], eigenVals[maxIdx] = eigenVals[maxIdx], eigenVals[i]
			for j := 0; j < m; j++ {
				tmpVecs[j][i], tmpVecs[j][maxIdx] = tmpVecs[j][maxIdx], tmpVecs[j][i]
			}
		}
	}
	eigenVecs := make([][1536]float32, m)
	for i := 0; i < m; i++ {
		eigenVecs[i] = [1536]float32(make([]float32, m))
		for j := 0; j < m; j++ {
			eigenVecs[i][j] = tmpVecs[j][i]
		}
	}

	// Project the input embeddings onto the first three principal components
	dim := 3
	proj := make([][]float32, n)
	for i := 0; i < n; i++ {
		proj[i] = make([]float32, dim)
		for j := 0; j < dim; j++ {
			for k := 0; k < m; k++ {
				proj[i][j] += (embeddings[i][k] - means[k]) * eigenVecs[j][k]
			}
		}
	}

	return proj
}
