//go:generate go run -tags=dev config_generate.go
package cli

const defaultCfg = `# Configuration for martian, tool for Stationeers localization.
# See https://github.com/st-l10n/martian for reference.

# List of available languages.
languages:
  - code: RU
    name: Russian
    font: russian

  - code: EN
    name: English
    font: english

  - code: FR
    name: French
    font: extended

  - code: DE
    name: German
    font: extended

  - code: IT
    name: Italian
    font: extended

  - code: KN
    name: Japanese
    font: cjk
    locale: ja

  - code: KO
    name: Korean
    font: hangul

  - code: PL
    name: Polish
    font: extended

  - code: PT
    name: Portuges
    prefix: portuguese
    locale: pt-PT

  - code: CN
    name: Simplified Chinese
    font: cjk
    locale: zh-CN

  - code: TW
    name: Traditional Chinese
    prefix: traditional-chinese
    font: cjk
    locale: zh-TW

  # Finnish
  - code: FI
    name: Suomi

  - code: SK
    name: Slovak
    font: russian

  - code: CS
    name: Czech
    font: extended
  
  - code: ES
    name: Spanish
    font: extended

# Simplified translation parts.
# If simplified, the "msgid" value of translation is set to simplified
# relative path of translated element, not the original text.
#
# Like for Reagents/RecordReagent with Key=Flour, the id for Unit will be "Flour.Unit"
# instead of "g".
#
# The "Tips" part is always assumed as non-simplified.
simplified:
  - Reagents.Unit
  - Keys

`
