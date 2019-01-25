/*
Package defaulter helps to set default.
It recursively walk through any types and try to apply defaults.

A default is applied if the type implement the setDefaultFunc interface
and the actual value is the zero value of the type.

Applying defaults

Considering the following structure:

	type Furniture struct{
		Color       Color
		Weight      float64
		IsAvailable bool
	}

	func (f *Furniture) SetDefault() {
		f.IsAvailable = true
	}


	type (
		Color string

		ColorUnknown = Color("unknown")
		ColorBlack   = Color("black")
		ColorGreen   = Color("green")
	)

	func (c *Color) SetDefault() {
		*c = ColorUnknown
	}

If we want to apply defaults to an instance of furniture, we just have to:

	var f Furniture
	defaulter.SetDefault(&f)
	// f has an unknown color and is available

Special cases

We already see that the SetDefault method will be call
if any value is the zero value and value's type implements SetDefault.
There is two exception:

If the type is a pointer:

	type Furniture struct{
		IsAvailable  bool
		OwnerAddress *Address
	}

	type Address string

	func (a *Address) SetDefault() {
		*a = "nowhere"
	}

	func foo() {
		var f Funiture

		defaulter.SetDefault(&f)
		// even if Address type implements SetDefault, the method will not be call
		// as f.OwnerAddress is nil, and we may want to keep it that way
	}


If it's a struct field with the `default:"-"` tag:

	type Furniture struct{
		IsAvailable  bool
		OwnerAddress Address `default:"-"`
	}

	func (a *Address) SetDefault() {
		*a = "nowhere"
	}

	func (f *Furniture) SetDefault() {
		f.IsAvailable = true
	}

	funf foo() {
		var f Furniture

		defaulter.SetDefault(&f)
		// here f.IsAvailable is true
		// but f.OwnderAddress still contains the zero value, because of the `nodefault` tag
	}
*/
package defaulter
