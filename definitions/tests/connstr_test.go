package main

//func TestFromString(t *testing.T) {
//	cases := []struct {
//		name    string
//		connStr string
//		pairs   []types.Pair
//		err     error
//	}{
//		{
//			"empty",
//			"",
//			nil,
//			services.ErrConnectionStringInvalid,
//		},
//		{
//			"simplest",
//			"tests://",
//			nil,
//			nil,
//		},
//		{
//			"only options",
//			"tests://?size=200",
//			[]types.Pair{
//				pairs.WithSize(200),
//			},
//			nil,
//		},
//		{
//			"only root dir",
//			"tests:///",
//			[]types.Pair{
//				pairs.WithWorkDir("/"),
//			},
//			nil,
//		},
//		{
//			"end with ?",
//			"tests:///?",
//			[]types.Pair{
//				pairs.WithWorkDir("/"),
//			},
//			nil,
//		},
//		{
//			"stupid, but valid (ignored)",
//			"tests:///?&&&",
//			[]types.Pair{
//				pairs.WithWorkDir("/"),
//			},
//			nil,
//		},
//		{
//			"value can contain all characters except &",
//			"tests:///?string_pair=a=b:/c?d&size=200",
//			[]types.Pair{
//				pairs.WithWorkDir("/"),
//				WithStringPair("a=b:/c?d"),
//				pairs.WithSize(200),
//			},
//			nil,
//		},
//		{
//			"full format",
//			"tests://abc/tmp/tmp1?size=200&storage_class=sc",
//			[]types.Pair{
//				pairs.WithName("abc"),
//				pairs.WithWorkDir("/tmp/tmp1"),
//				pairs.WithSize(200),
//				WithStorageClass("sc"),
//			},
//			nil,
//		},
//		{
//			"duplicate key, appear in order (finally, first will be picked)",
//			"tests://abc/tmp/tmp1?size=200&name=def&size=300",
//			[]types.Pair{
//				pairs.WithName("abc"),
//				pairs.WithWorkDir("/tmp/tmp1"),
//				pairs.WithSize(200),
//				pairs.WithName("def"),
//				pairs.WithSize(300),
//			},
//			nil,
//		},
//		{
//			"not registered pair",
//			"tests://abc/tmp?not_a_pair=a",
//			nil,
//			services.ErrConnectionStringInvalid,
//		},
//		{
//			"key without value (not registered pair)",
//			"tests://abc/tmp?not_a_pair&&",
//			nil,
//			services.ErrConnectionStringInvalid,
//		},
//		{
//			"not parsable pair",
//			"tests://abc/tmp?io_call_back=a",
//			nil,
//			services.ErrConnectionStringInvalid,
//		},
//		{
//			"key with features",
//			"tests://abc/tmp?enable_loose_pair",
//			[]types.Pair{
//				pairs.WithName("abc"),
//				pairs.WithWorkDir("/tmp"),
//				WithEnableLoosePair(),
//			},
//			nil,
//		},
//		{
//			"key with default paris",
//			"tests://abc/tmp?default_storage_class=STANDARD",
//			[]types.Pair{
//				pairs.WithName("abc"),
//				pairs.WithWorkDir("/tmp"),
//				WithDefaultStorageClass("STANDARD"),
//			},
//			nil,
//		},
//	}
//
//	for _, tt := range cases {
//		t.Run(tt.name, func(t *testing.T) {
//			t.Run("NewServicerFromString", func(t *testing.T) {
//				servicer, err := services.NewServicerFromString(tt.connStr)
//				service, ok := servicer.(*Service)
//
//				if tt.err == nil {
//					assert.Nil(t, err)
//					assert.True(t, ok)
//				} else {
//					assert.True(t, errors.Is(err, tt.err))
//					return
//				}
//
//				assert.Equal(t, service.Pairs, tt.pairs)
//			})
//			t.Run("NewStoragerFromString", func(t *testing.T) {
//				storager, err := services.NewStoragerFromString(tt.connStr)
//				storage, ok := storager.(*Storage)
//
//				if tt.err == nil {
//					assert.Nil(t, err)
//					assert.True(t, ok)
//				} else {
//					assert.True(t, errors.Is(err, tt.err))
//					return
//				}
//
//				assert.Equal(t, storage.Pairs, tt.pairs)
//			})
//
//		})
//	}
//}
