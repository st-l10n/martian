# Stationeers localication tools

Tools that prepare localization resources.

## Usage
```bash
# install
$ go get github.com/st-10n/martian
# prepare environment
$ mkdir locales
$ export INPUT_DIR=~/.local/share/Steam/steamapps/common/Stationeers/rocketstation_Data/StreamingAssets/Language/
# generate locales (--limit is optional param)
$ martian gen --input $INPUT_DIR -o locales --limit en,sk,ru,de
postfixes: [.xml _keys.xml _tips.xml]
Language: Russian
  prefix: russian
  code: RU
  locale: ru
  entries: 1641
Language: English
  prefix: english
  code: EN
  locale: en
  entries: 1651
Language: German
  prefix: german
  code: DE
  locale: de
  entries: 1553
Language: Slovak
  prefix: slovak
  code: SK
  locale: sk
  entries: 1266
```