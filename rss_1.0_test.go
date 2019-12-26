package rss

import (
	"io/ioutil"
	"testing"
)

type rdfResult struct {
	Title string
	Link  string
}

func TestParseRDFRSS10(t *testing.T) {
	res := getRdfTestResults()

	data, err := ioutil.ReadFile("testdata/rss_rdf")
	if err != nil {
		t.Fatalf("Reading test file %v", err)
	}

	feed, err := Parse(data, nil, "")
	if err != nil {
		t.Error("Should not error on parsing", err)
	}

	for index, expectedItem := range res {
		title := feed.Items[index].Title
		link := feed.Items[index].Link
		id := feed.Items[index].ID
		if title != expectedItem.Title {
			t.Errorf("[title] Expected %s got %s", expectedItem.Title, title)
		}
		if link != expectedItem.Link {
			t.Errorf("[link] Expected %s got %s", expectedItem.Link, link)
		}
		if id != expectedItem.Link {
			t.Errorf("[id] Expected %s got %s", expectedItem.Link, id)
		}
	}
}

func getRdfTestResults() []rdfResult {
	return []rdfResult{
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1609976?af=R",
			Title: "An Expectation Conditional Maximization Approach for Gaussian Graphical Models",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1609977?af=R",
			Title: "Beyond Prediction: A Framework for Inference With Variational Approximations in Mixture Models",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1598872?af=R",
			Title: "Adaptive Incremental Mixture Markov Chain Monte Carlo",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1592756?af=R",
			Title: "Incremental Mixture Importance Sampling With Shotgun Optimization",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1592757?af=R",
			Title: "Easily Parallelizable and Distributable Class of Algorithms for Structured Sparsity, with Optimal Acceleration",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1594835?af=R",
			Title: "Damped Anderson Acceleration With Restarts and Monotonicity Control for Accelerating EM and EM-like Algorithms",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1598871?af=R",
			Title: "Projection Pursuit Based on Gaussian Mixtures and Evolutionary Algorithms",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1601097?af=R",
			Title: "A Metaheuristic Adaptive Cubature Based Algorithm to Find Bayesian Optimal Designs for Nonlinear Models",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1598869?af=R",
			Title: "Influence Diagnostics for High-Dimensional Lasso Regression",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1586714?af=R",
			Title: "Distributed Generalized Cross-Validation for Divide-and-Conquer Kernel Ridge Regression and Its Asymptotic Optimality",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1598870?af=R",
			Title: "Component-Based Regularization of Multivariate Generalized Linear Mixed Models",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1593179?af=R",
			Title: "Simultaneous Variable and Covariance Selection With the Multivariate Spike-and-Slab LASSO",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1604374?af=R",
			Title: "The Generalized Ridge Estimator of the Inverse Covariance Matrix",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1607744?af=R",
			Title: "Simultaneous Registration and Clustering for Multidimensional Functional Data",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1592758?af=R",
			Title: "Flexible and Interpretable Models for Survival Data",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1594834?af=R",
			Title: "Stable Multiple Time Step Simulation/Prediction From Lagged Dynamic Network Regression Models",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1593180?af=R",
			Title: "Improving Spectral Clustering Using the Asymptotic Value of the Normalized Cut",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1604373?af=R",
			Title: "Variable-Domain Functional Principal Component Analysis",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1610312?af=R",
			Title: "Fast Generalized Linear Models by Database Sampling and One-Step Polishing",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1637746?af=R",
			Title: "Good Plot Symbols by Default",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1631620?af=R",
			Title: "Correction",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1654746?af=R",
			Title: "Correction",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1695454?af=R",
			Title: "Correction",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1694523?af=R",
			Title: "Inferring Influence Networks From Longitudinal Bipartite Relational Data",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1694524?af=R",
			Title: "First- and Second-Order Characteristics of Spatio-Temporal Point Processes on Linear Networks",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1686988?af=R",
			Title: "A Bayesian Time-Varying Coefficient Model for Multitype Recurrent Events",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1677474?af=R",
			Title: "Testing One Hypothesis Multiple Times: The Multidimensional Case",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1677243?af=R",
			Title: "Heteroscedastic BART via Multiplicative Regression Trees",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1677244?af=R",
			Title: "Longitudinal Principal Component Analysis With an Application to Marketing Data",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1665537?af=R",
			Title: "The Rational SPDE Approach for Gaussian Random Fields With General Smoothness",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1665535?af=R",
			Title: "A Logistic Factorization Model for Recommender Systems With Multinomial Responses",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1668800?af=R",
			Title: "Simulating Markov Random Fields With a Conclique-Based Gibbs Sampler",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1665534?af=R",
			Title: "Fast Nonseparable Gaussian Stochastic Process With Application to Methylation Level Interpolation",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1665536?af=R",
			Title: "Consistent Blind Image Deblurring Using Jump-Preserving Extrapolation",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1660180?af=R",
			Title: "Valid Inference Corrected for Outlier Removal",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1647846?af=R",
			Title: "Estimating the Number of Clusters Using Cross-Validation",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1660178?af=R",
			Title: "Solving Fused Penalty Estimation Problems via Block Splitting Algorithms",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1660179?af=R",
			Title: "Compressed and Penalized Linear Regression",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1654881?af=R",
			Title: "Area-Proportional Visualization for Circular Data",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1654880?af=R",
			Title: "Scalable Bayesian Inference for Coupled Hidden Markov and Semi-Markov Models",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1647216?af=R",
			Title: "Parallelization of a Common Changepoint Detection Method",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1647849?af=R",
			Title: "Bivariate Residual Plots With Simulation Polygons",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1652616?af=R",
			Title: "Scalable Gaussian Process Computations Using Hierarchical Matrices",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1648271?af=R",
			Title: "A Pliable Lasso",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1647215?af=R",
			Title: "Efficient Construction of Test Inversion Confidence Intervals Using Quantile Regression",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1647848?af=R",
			Title: "Estimating Time-Varying Graphical Models",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1647847?af=R",
			Title: "Bayesian Model Averaging Over Tree-based Dependence Structures for Multivariate Extremes",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1637749?af=R",
			Title: "Testing Sparsity-Inducing Penalties",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1637747?af=R",
			Title: "Bayesian Deep Net GLM and GLMM",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1637748?af=R",
			Title: "Diagonal Discriminant Analysis With Feature Selection for High-Dimensional Data",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1629941?af=R",
			Title: "A Function Emulation Approach for Doubly Intractable Distributions",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1624365?af=R",
			Title: "BIVAS: A Scalable Bayesian Method for Bi-Level Variable Selection With Applications",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1624366?af=R",
			Title: "Scalable Bayesian Nonparametric Clustering and Classification",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1629942?af=R",
			Title: "Scalable Visualization Methods for Modern Generalized Additive Models",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1629943?af=R",
			Title: "Dynamic Visualization and Fast Computation for Convex Clustering via Algorithmic Regularization",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1624294?af=R",
			Title: "Scalable Bayesian Regression in High Dimensions With Multiple Data Sources",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1617159?af=R",
			Title: "A Semiparametric Bayesian Approach to Dropout in Longitudinal Studies With Auxiliary Covariates",
		},
		rdfResult{
			Link:  "https://www.tandfonline.com/doi/full/10.1080/10618600.2019.1617160?af=R",
			Title: "Anomaly Detection in Streaming Nonstationary Temporal Data",
		},
	}
}
