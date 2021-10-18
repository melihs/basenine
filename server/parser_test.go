package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicBoolean(t *testing.T) {
	text := `
http or !amqp
	`
	expr := Parse(text)
	// repr.Println(expr)

	val1 := "http"
	val2 := "amqp"

	expect := &Expression{
		Equality: &Equality{
			Comparison: &Comparison{
				Logical: &Logical{
					Unary: &Unary{
						Primary: &Primary{
							CallExpression: &CallExpression{
								Identifier: &val1,
							},
						},
					},
					Op: "or",
					Next: &Logical{
						Unary: &Unary{
							Op: "!",
							Unary: &Unary{
								Primary: &Primary{
									CallExpression: &CallExpression{
										Identifier: &val2,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expect, expr)
}

func TestBooleanLiterals(t *testing.T) {
	text := `
true and false
	`
	expr := Parse(text)
	// repr.Println(expr)

	val1 := true

	expect := &Expression{
		Equality: &Equality{
			Comparison: &Comparison{
				Logical: &Logical{
					Unary: &Unary{
						Primary: &Primary{
							Bool: &val1,
						},
					},
					Op: "and",
					Next: &Logical{
						Unary: &Unary{
							Primary: &Primary{},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expect, expr)
}

func TestCompoundBoolean(t *testing.T) {
	text := `
true and 5 == a
	`
	expr := Parse(text)
	// repr.Println(expr)

	val1 := true
	val2 := float64(5)
	val3 := "a"

	expect := &Expression{
		Equality: &Equality{
			Comparison: &Comparison{
				Logical: &Logical{
					Unary: &Unary{
						Primary: &Primary{
							Bool: &val1,
						},
					},
					Op: "and",
					Next: &Logical{
						Unary: &Unary{
							Primary: &Primary{
								Number: &val2,
							},
						},
					},
				},
			},
			Op: "==",
			Next: &Equality{
				Comparison: &Comparison{
					Logical: &Logical{
						Unary: &Unary{
							Primary: &Primary{
								CallExpression: &CallExpression{
									Identifier: &val3,
								},
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expect, expr)
}

func TestNegatedCompoundBoolean(t *testing.T) {
	text := `
true and !(5 == a)
	`
	expr := Parse(text)
	// repr.Println(expr)

	val1 := true
	val2 := float64(5)
	val3 := "a"

	expect := &Expression{
		Equality: &Equality{
			Comparison: &Comparison{
				Logical: &Logical{
					Unary: &Unary{
						Primary: &Primary{
							Bool: &val1,
						},
					},
					Op: "and",
					Next: &Logical{
						Unary: &Unary{
							Op: "!",
							Unary: &Unary{
								Primary: &Primary{
									SubExpression: &Expression{
										Equality: &Equality{
											Comparison: &Comparison{
												Logical: &Logical{
													Unary: &Unary{
														Primary: &Primary{
															Number: &val2,
														},
													},
												},
											},
											Op: "==",
											Next: &Equality{
												Comparison: &Comparison{
													Logical: &Logical{
														Unary: &Unary{
															Primary: &Primary{
																CallExpression: &CallExpression{
																	Identifier: &val3,
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expect, expr)
}

func TestSubExpression(t *testing.T) {
	text := `
	(a.b == "hello") and (x.y > 3.14)
	`
	expr := Parse(text)
	// repr.Println(expr)

	val1 := "a.b"
	val2 := "\"hello\""
	val3 := "x.y"
	val4 := 3.14

	expect := &Expression{
		Equality: &Equality{
			Comparison: &Comparison{
				Logical: &Logical{
					Unary: &Unary{
						Primary: &Primary{
							SubExpression: &Expression{
								Equality: &Equality{
									Comparison: &Comparison{
										Logical: &Logical{
											Unary: &Unary{
												Primary: &Primary{
													CallExpression: &CallExpression{
														Identifier: &val1,
													},
												},
											},
										},
									},
									Op: "==",
									Next: &Equality{
										Comparison: &Comparison{
											Logical: &Logical{
												Unary: &Unary{
													Primary: &Primary{
														String: &val2,
													},
												},
											},
										},
									},
								},
							},
						},
					},
					Op: "and",
					Next: &Logical{
						Unary: &Unary{
							Primary: &Primary{
								SubExpression: &Expression{
									Equality: &Equality{
										Comparison: &Comparison{
											Logical: &Logical{
												Unary: &Unary{
													Primary: &Primary{
														CallExpression: &CallExpression{
															Identifier: &val3,
														},
													},
												},
											},
											Op: ">",
											Next: &Comparison{
												Logical: &Logical{
													Unary: &Unary{
														Primary: &Primary{
															Number: &val4,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expect, expr)
}

func TestRegexLiteral(t *testing.T) {
	text := `
http.request == r"hello"
	`
	expr := Parse(text)
	// repr.Println(expr)

	val1 := "http.request"
	val2 := "\"hello\""

	expect := &Expression{
		Equality: &Equality{
			Comparison: &Comparison{
				Logical: &Logical{
					Unary: &Unary{
						Primary: &Primary{
							CallExpression: &CallExpression{
								Identifier: &val1,
							},
						},
					},
				},
			},
			Op: "==",
			Next: &Equality{
				Comparison: &Comparison{
					Logical: &Logical{
						Unary: &Unary{
							Primary: &Primary{
								Regex: &val2,
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expect, expr)
}

func TestComplexQuery(t *testing.T) {
	text := `
http and request.method == "GET" and request.path == "/example" and (request.query.a == "b" or request.headers.x == "y")
	`
	expr := Parse(text)
	// repr.Println(expr)

	val1 := "http"
	val2 := "request.method"
	val3 := "\"GET\""
	val4 := "request.path"
	val5 := "\"/example\""
	val6 := "request.query.a"
	val7 := "\"b\""
	val8 := "request.headers.x"
	val9 := "\"y\""

	expect := &Expression{
		Equality: &Equality{
			Comparison: &Comparison{
				Logical: &Logical{
					Unary: &Unary{
						Primary: &Primary{
							CallExpression: &CallExpression{
								Identifier: &val1,
							},
						},
					},
					Op: "and",
					Next: &Logical{
						Unary: &Unary{
							Primary: &Primary{
								CallExpression: &CallExpression{
									Identifier: &val2,
								},
							},
						},
					},
				},
			},
			Op: "==",
			Next: &Equality{
				Comparison: &Comparison{
					Logical: &Logical{
						Unary: &Unary{
							Primary: &Primary{
								String: &val3,
							},
						},
						Op: "and",
						Next: &Logical{
							Unary: &Unary{
								Primary: &Primary{
									CallExpression: &CallExpression{
										Identifier: &val4,
									},
								},
							},
						},
					},
				},
				Op: "==",
				Next: &Equality{
					Comparison: &Comparison{
						Logical: &Logical{
							Unary: &Unary{
								Primary: &Primary{
									String: &val5,
								},
							},
							Op: "and",
							Next: &Logical{
								Unary: &Unary{
									Primary: &Primary{
										SubExpression: &Expression{
											Equality: &Equality{
												Comparison: &Comparison{
													Logical: &Logical{
														Unary: &Unary{
															Primary: &Primary{
																CallExpression: &CallExpression{
																	Identifier: &val6,
																},
															},
														},
													},
												},
												Op: "==",
												Next: &Equality{
													Comparison: &Comparison{
														Logical: &Logical{
															Unary: &Unary{
																Primary: &Primary{
																	String: &val7,
																},
															},
															Op: "or",
															Next: &Logical{
																Unary: &Unary{
																	Primary: &Primary{
																		CallExpression: &CallExpression{
																			Identifier: &val8,
																		},
																	},
																},
															},
														},
													},
													Op: "==",
													Next: &Equality{
														Comparison: &Comparison{
															Logical: &Logical{
																Unary: &Unary{
																	Primary: &Primary{
																		String: &val9,
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expect, expr)
}

func TestSelectExpressionIndex(t *testing.T) {
	text := `
http.request.path[1] == "hello"
	`
	expr := Parse(text)
	// repr.Println(expr)

	val1 := "http.request.path"
	val2 := "\"hello\""

	expect := &Expression{
		Equality: &Equality{
			Comparison: &Comparison{
				Logical: &Logical{
					Unary: &Unary{
						Primary: &Primary{
							CallExpression: &CallExpression{
								Identifier: &val1,
								SelectExpression: &SelectExpression{
									Index: 1,
								},
							},
						},
					},
				},
			},
			Op: "==",
			Next: &Equality{
				Comparison: &Comparison{
					Logical: &Logical{
						Unary: &Unary{
							Primary: &Primary{
								String: &val2,
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expect, expr)
}

func TestSelectExpressionKey(t *testing.T) {
	text := `
!http.request.headers["user-agent"] == "kube-probe"
	`
	expr := Parse(text)
	// repr.Println(expr)

	val1 := "http.request.headers"
	val2 := "\"user-agent\""
	val3 := "\"kube-probe\""

	expect := &Expression{
		Equality: &Equality{
			Comparison: &Comparison{
				Logical: &Logical{
					Unary: &Unary{
						Op: "!",
						Unary: &Unary{
							Primary: &Primary{
								CallExpression: &CallExpression{
									Identifier: &val1,
									SelectExpression: &SelectExpression{
										Key: &val2,
									},
								},
							},
						},
					},
				},
			},
			Op: "==",
			Next: &Equality{
				Comparison: &Comparison{
					Logical: &Logical{
						Unary: &Unary{
							Primary: &Primary{
								String: &val3,
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expect, expr)
}

func TestFunctionCall(t *testing.T) {
	text := `
a.b(3, 5)
	`
	expr := Parse(text)
	// repr.Println(expr)

	val1 := "a.b"
	val2 := float64(3)
	val3 := float64(5)

	expect := &Expression{
		Equality: &Equality{
			Comparison: &Comparison{
				Logical: &Logical{
					Unary: &Unary{
						Primary: &Primary{
							CallExpression: &CallExpression{
								Identifier: &val1,
								Parameters: []*Parameter{
									{
										Expression: &Expression{
											Equality: &Equality{
												Comparison: &Comparison{
													Logical: &Logical{
														Unary: &Unary{
															Primary: &Primary{
																Number: &val2,
															},
														},
													},
												},
											},
										},
									},
									{
										Expression: &Expression{
											Equality: &Equality{
												Comparison: &Comparison{
													Logical: &Logical{
														Unary: &Unary{
															Primary: &Primary{
																Number: &val3,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expect, expr)
}

func TestSelectExpressionChainFunction(t *testing.T) {
	text := `
!http or !http.request.headers["user-agent"].startsWith("kube-probe")
	`
	expr := Parse(text)
	// repr.Println(expr)

	val1 := "http"
	val2 := "http.request.headers"
	val3 := "\"user-agent\""
	val4 := "startsWith"
	val5 := "\"kube-probe\""

	expect := &Expression{
		Equality: &Equality{
			Comparison: &Comparison{
				Logical: &Logical{
					Unary: &Unary{
						Op: "!",
						Unary: &Unary{
							Primary: &Primary{
								CallExpression: &CallExpression{
									Identifier: &val1,
								},
							},
						},
					},
					Op: "or",
					Next: &Logical{
						Unary: &Unary{
							Op: "!",
							Unary: &Unary{
								Primary: &Primary{
									CallExpression: &CallExpression{
										Identifier: &val2,
										SelectExpression: &SelectExpression{
											Key: &val3,
											Expression: &Expression{
												Equality: &Equality{
													Comparison: &Comparison{
														Logical: &Logical{
															Unary: &Unary{
																Primary: &Primary{
																	CallExpression: &CallExpression{
																		Identifier: &val4,
																		Parameters: []*Parameter{
																			{
																				Expression: &Expression{
																					Equality: &Equality{
																						Comparison: &Comparison{
																							Logical: &Logical{
																								Unary: &Unary{
																									Primary: &Primary{
																										String: &val5,
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expect, expr)
}

func TestRulesAssertionSyntax(t *testing.T) {
	text := `
rule(
	description: "Holy in name property",
	query: http and service == r"catalogue.*" and request.path == r"catalogue.*" and response.headers["content-type"].contains("application/json"),
	assert: response.body.name == "Holy"
)
and
rule(
	description: "Content Length header",
	query: http,
	assert: response.headers["content-length"] == r"(\\d+(?:\\.\\d+)?)"
)
and
rule(
	description: "Latency test",
	query: http and service == r"carts.*",
	assert: response.elapsedTime >= 1
)
	`
	expr := Parse(text)
	// repr.Println(expr)

	val1 := "rule"
	val2 := "description"
	val3 := "\"Holy in name property\""
	val4 := "query"
	val5 := "http"
	val6 := "service"
	val7 := "\"catalogue.*\""
	val8 := "request.path"
	val9 := "response.headers"
	val10 := "\"content-type\""
	val11 := "contains"
	val12 := "\"application/json\""
	val13 := "assert"
	val14 := "response.body.name"
	val15 := "\"Holy\""
	val16 := "\"Content Length header\""
	val17 := "\"content-length\""
	val18 := "\"(\\\\d+(?:\\\\.\\\\d+)?)\""
	val19 := "\"Latency test\""
	val20 := "service"
	val21 := "\"carts.*\""
	val22 := "response.elapsedTime"
	val23 := float64(1)

	expect := &Expression{
		Equality: &Equality{
			Comparison: &Comparison{
				Logical: &Logical{
					Unary: &Unary{
						Primary: &Primary{
							CallExpression: &CallExpression{
								Identifier: &val1,
								Parameters: []*Parameter{
									{
										Tag: &val2,
										Expression: &Expression{
											Equality: &Equality{
												Comparison: &Comparison{
													Logical: &Logical{
														Unary: &Unary{
															Primary: &Primary{
																String: &val3,
															},
														},
													},
												},
											},
										},
									},
									{
										Tag: &val4,
										Expression: &Expression{
											Equality: &Equality{
												Comparison: &Comparison{
													Logical: &Logical{
														Unary: &Unary{
															Primary: &Primary{
																CallExpression: &CallExpression{
																	Identifier: &val5,
																},
															},
														},
														Op: "and",
														Next: &Logical{
															Unary: &Unary{
																Primary: &Primary{
																	CallExpression: &CallExpression{
																		Identifier: &val6,
																	},
																},
															},
														},
													},
												},
												Op: "==",
												Next: &Equality{
													Comparison: &Comparison{
														Logical: &Logical{
															Unary: &Unary{
																Primary: &Primary{
																	Regex: &val7,
																},
															},
															Op: "and",
															Next: &Logical{
																Unary: &Unary{
																	Primary: &Primary{
																		CallExpression: &CallExpression{
																			Identifier: &val8,
																		},
																	},
																},
															},
														},
													},
													Op: "==",
													Next: &Equality{
														Comparison: &Comparison{
															Logical: &Logical{
																Unary: &Unary{
																	Primary: &Primary{
																		Regex: &val7,
																	},
																},
																Op: "and",
																Next: &Logical{
																	Unary: &Unary{
																		Primary: &Primary{
																			CallExpression: &CallExpression{
																				Identifier: &val9,
																				SelectExpression: &SelectExpression{
																					Key: &val10,
																					Expression: &Expression{
																						Equality: &Equality{
																							Comparison: &Comparison{
																								Logical: &Logical{
																									Unary: &Unary{
																										Primary: &Primary{
																											CallExpression: &CallExpression{
																												Identifier: &val11,
																												Parameters: []*Parameter{
																													{
																														Expression: &Expression{
																															Equality: &Equality{
																																Comparison: &Comparison{
																																	Logical: &Logical{
																																		Unary: &Unary{
																																			Primary: &Primary{
																																				String: &val12,
																																			},
																																		},
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									{
										Tag: &val13,
										Expression: &Expression{
											Equality: &Equality{
												Comparison: &Comparison{
													Logical: &Logical{
														Unary: &Unary{
															Primary: &Primary{
																CallExpression: &CallExpression{
																	Identifier: &val14,
																},
															},
														},
													},
												},
												Op: "==",
												Next: &Equality{
													Comparison: &Comparison{
														Logical: &Logical{
															Unary: &Unary{
																Primary: &Primary{
																	String: &val15,
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					Op: "and",
					Next: &Logical{
						Unary: &Unary{
							Primary: &Primary{
								CallExpression: &CallExpression{
									Identifier: &val1,
									Parameters: []*Parameter{
										{
											Tag: &val2,
											Expression: &Expression{
												Equality: &Equality{
													Comparison: &Comparison{
														Logical: &Logical{
															Unary: &Unary{
																Primary: &Primary{
																	String: &val16,
																},
															},
														},
													},
												},
											},
										},
										{
											Tag: &val4,
											Expression: &Expression{
												Equality: &Equality{
													Comparison: &Comparison{
														Logical: &Logical{
															Unary: &Unary{
																Primary: &Primary{
																	CallExpression: &CallExpression{
																		Identifier: &val5,
																	},
																},
															},
														},
													},
												},
											},
										},
										{
											Tag: &val13,
											Expression: &Expression{
												Equality: &Equality{
													Comparison: &Comparison{
														Logical: &Logical{
															Unary: &Unary{
																Primary: &Primary{
																	CallExpression: &CallExpression{
																		Identifier: &val9,
																		SelectExpression: &SelectExpression{
																			Key: &val17,
																		},
																	},
																},
															},
														},
													},
													Op: "==",
													Next: &Equality{
														Comparison: &Comparison{
															Logical: &Logical{
																Unary: &Unary{
																	Primary: &Primary{
																		Regex: &val18,
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						Op: "and",
						Next: &Logical{
							Unary: &Unary{
								Primary: &Primary{
									CallExpression: &CallExpression{
										Identifier: &val1,
										Parameters: []*Parameter{
											{
												Tag: &val2,
												Expression: &Expression{
													Equality: &Equality{
														Comparison: &Comparison{
															Logical: &Logical{
																Unary: &Unary{
																	Primary: &Primary{
																		String: &val19,
																	},
																},
															},
														},
													},
												},
											},
											{
												Tag: &val4,
												Expression: &Expression{
													Equality: &Equality{
														Comparison: &Comparison{
															Logical: &Logical{
																Unary: &Unary{
																	Primary: &Primary{
																		CallExpression: &CallExpression{
																			Identifier: &val5,
																		},
																	},
																},
																Op: "and",
																Next: &Logical{
																	Unary: &Unary{
																		Primary: &Primary{
																			CallExpression: &CallExpression{
																				Identifier: &val20,
																			},
																		},
																	},
																},
															},
														},
														Op: "==",
														Next: &Equality{
															Comparison: &Comparison{
																Logical: &Logical{
																	Unary: &Unary{
																		Primary: &Primary{
																			Regex: &val21,
																		},
																	},
																},
															},
														},
													},
												},
											},
											{
												Tag: &val13,
												Expression: &Expression{
													Equality: &Equality{
														Comparison: &Comparison{
															Logical: &Logical{
																Unary: &Unary{
																	Primary: &Primary{
																		CallExpression: &CallExpression{
																			Identifier: &val22,
																		},
																	},
																},
															},
															Op: ">=",
															Next: &Comparison{
																Logical: &Logical{
																	Unary: &Unary{
																		Primary: &Primary{
																			Number: &val23,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expect, expr)
}
