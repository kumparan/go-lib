package token

func getRole(level int) string {
	var roleLevel =  map[string][]int{
		"basic": {0, 19},
		"superadmin": {1},
		"admin": {2,3,12,13},
		"contentIntelligent": {3,4,5},
		"editorial": {6,7,8,9,10,11,17,18},
		"communityEngagement": {12, 13, 14, 15, 16},
	}

	role, found := mapKey(roleLevel,level)

  if found == false {
    role = "basic";
	}

	return role
}

func mapKey(m map[string][]int, value int) (key string, ok bool) {
  for k, v := range m {
		for _, v2 := range v {
			if v2 == value { 
				key = k
				ok = true
				return
			}
		}
  }
  return
}