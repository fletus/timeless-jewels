package main

// The tests and benchmarks in this package assert against the SHIPPED tables, so they need the
// embedded assets installed into the data package.
//
// ⚠️ THIS USED TO HAPPEN BY ITSELF. Before the data package was split, importing it ran an
// init() that decompressed the assets, so every consumer got the tables whether it wanted them
// or not. Now the assets live in data/embedded and this import is what installs them. Deleting
// it does not break the build — it makes every lookup return nil, which the value assertions
// below catch, but only loudly because they assert values.
import _ "github.com/Vilsol/timeless-jewels/data/embedded"
