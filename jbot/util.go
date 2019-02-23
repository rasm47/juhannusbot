package jbot

// stringSlicesAreEqual returns true if two string slices (a and b) are fully equal.
func stringSlicesAreEqual(a []string, b []string) bool {
    if (a == nil) != (b == nil) { 
        return false; 
    }

    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}
