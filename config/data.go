package config

var OptionsMap = map[string]string{
	".flake8":        FlakeData,
	"License":        LicenseData,
	"pyproject.toml": PyProjectTomlData,
}

var FlakeData = `[flake8]
ignore = D203
exclude = .git,__pycache__,docs/source/conf.py,old,build,dist,venv,migrations,deploy
max-complexity = 10
max-line-length = 160
`

var LicenseData = `The MIT License (MIT)
Copyright (c) 2020, DOE Data Science

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
`

var PyProjectTomlData = `[tool.black]
line-length = 160
target-version = ['py38']
include = '\.pyi?$'
extend-exclude = '''
^/venv
'''
`
