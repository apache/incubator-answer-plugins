/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package lark

import "encoding/json"

// CardLink represents the URLs for different platforms.
type CardLink struct {
	URL        string `json:"url,omitempty"`
	PCURL      string `json:"pc_url,omitempty"`
	IOSURL     string `json:"ios_url,omitempty"`
	AndroidURL string `json:"android_url,omitempty"`
}

// Text represents a text component with optional internationalization.
type Text struct {
	Tag       string `json:"tag,omitempty"`
	Content   string `json:"content,omitempty"`
	I18n      *I18n  `json:"i18n,omitempty"`
	TextSize  string `json:"text_size,omitempty"`
	TextAlign string `json:"text_align,omitempty"`
	TextColor string `json:"text_color,omitempty"`
	Icon      *Icon  `json:"icon,omitempty"`
}

// I18n represents internationalized text.
type I18n struct {
	ZhCn string `json:"zh_cn,omitempty"`
	EnUs string `json:"en_us,omitempty"`
}

// TextTag represents a text tag component.
type TextTag struct {
	Tag   string `json:"tag,omitempty"`
	Text  *Text  `json:"text,omitempty"`
	Color string `json:"color,omitempty"`
}

// MarshalJSON customizes the JSON encoding for TextTag.
func (tt *TextTag) MarshalJSON() ([]byte, error) {
	tt.Tag = "text_tag"
	return json.Marshal(*tt)
}

// Icon represents an icon component.
type Icon struct {
	Tag    string `json:"tag,omitempty"`
	Token  string `json:"token,omitempty"`
	Color  string `json:"color,omitempty"`
	ImgKey string `json:"img_key,omitempty"`
	Style  *struct {
		Color Color `json:"color,omitempty"`
	} `json:"style,omitempty"`
}

// Header represents the header component of a card.
type Header struct {
	Title       *Text     `json:"title,omitempty"`
	Subtitle    *Text     `json:"subtitle,omitempty"`
	TextTagList []TextTag `json:"text_tag_list,omitempty"`
	Template    Template  `json:"template,omitempty"`
	UdIcon      *Icon     `json:"ud_icon,omitempty"`
}

// Element represents a generic element in a card.
type Element struct {
	*PlainText
	*Button
}

// MarshalJSON customizes the JSON encoding for Element.
func (e *Element) MarshalJSON() ([]byte, error) {
	if e.PlainText != nil {
		return json.Marshal(e.PlainText)
	}
	return json.Marshal(e.Button)
}

// Column represents a column in a ColumnSet.
type Column struct {
	Tag             string    `json:"tag,omitempty"`
	Elements        []Element `json:"elements,omitempty"`
	Width           string    `json:"width,omitempty"`
	Weight          int       `json:"weight,omitempty"`
	BackgroundStyle string    `json:"background_style,omitempty"`
	VerticalAlign   string    `json:"vertical_align,omitempty"`
	VerticalSpacing string    `json:"vertical_spacing,omitempty"`
	Padding         string    `json:"padding,omitempty"`
}

// MarshalJSON customizes the JSON encoding for Column.
func (c *Column) MarshalJSON() ([]byte, error) {
	c.Tag = "column"
	return json.Marshal(*c)
}

// ColumnSet represents a set of columns.
type ColumnSet struct {
	*Show
	*Action
}

// MarshalJSON customizes the JSON encoding for ColumnSet.
func (cs *ColumnSet) MarshalJSON() ([]byte, error) {
	if cs.Show != nil {
		cs.Show.Tag = "column_set"
		return json.Marshal(*cs.Show)
	}
	if cs.Action != nil {
		cs.Action.Tag = "action"
		return json.Marshal(*cs.Action)
	}
	return nil, nil
}

// Show represents the display properties of a ColumnSet.
type Show struct {
	Tag               string   `json:"tag"`
	FlexMode          string   `json:"flex_mode,omitempty"`
	HorizontalSpacing string   `json:"horizontal_spacing,omitempty"`
	BackgroundStyle   string   `json:"background_style,omitempty"`
	Columns           []Column `json:"columns,omitempty"`
}

// Action represents actions in a ColumnSet.
type Action struct {
	Tag     string    `json:"tag"`
	Actions []*Button `json:"actions,omitempty"`
}

// I18nElements represents internationalized elements.
type I18nElements struct {
	ZhCn []ColumnSet `json:"zh_cn,omitempty"`
	EnUs []ColumnSet `json:"en_us,omitempty"`
}

// Behavior represents the behavior of a button.
type Behavior struct {
	Type       string `json:"type"`
	DefaultURL string `json:"default_url"`
	AndroidURL string `json:"android_url"`
	IOSURL     string `json:"ios_url"`
	PCURL      string `json:"pc_url"`
}

// Button represents a button component.
type Button struct {
	Tag       string         `json:"tag,omitempty"`
	Width     string         `json:"width,omitempty"`
	Text      *Text          `json:"text,omitempty"`
	Behaviors []Behavior     `json:"behaviors,omitempty"`
	Type      string         `json:"type,omitempty"`
	HoverTips *Text          `json:"hover_tips,omitempty"`
	Value     map[string]any `json:"value,omitempty"`
}

// MarshalJSON customizes the JSON encoding for Button.
func (b *Button) MarshalJSON() ([]byte, error) {
	b.Tag = "button"
	return json.Marshal(*b)
}

// PlainText represents plain text component.
type PlainText struct {
	Tag  string `json:"tag,omitempty"`
	Text *Text  `json:"text,omitempty"`
	Icon *Icon  `json:"icon,omitempty"`
}

// Summary represents the summary information of a card.
type Summary struct {
	Content     string            `json:"content,omitempty"`
	I18nContent map[string]string `json:"i18n_content,omitempty"`
}

// TextSize represents the custom text size configuration.
type TextSize struct {
	Default string `json:"default,omitempty"`
	PC      string `json:"pc,omitempty"`
	Mobile  string `json:"mobile,omitempty"`
}

// ConfigColor represents the custom color configuration.
type ConfigColor struct {
	LightMode string `json:"light_mode,omitempty"`
	DarkMode  string `json:"dark_mode,omitempty"`
}

// Style represents the custom font size and color configuration.
type Style struct {
	TextSize map[string]TextSize    `json:"text_size,omitempty"`
	Color    map[string]ConfigColor `json:"color,omitempty"`
}

// Config represents the configuration of a card.
type Config struct {
	StreamingMode            *bool    `json:"streaming_mode,omitempty"`
	Summary                  *Summary `json:"summary,omitempty"`
	EnableForward            *bool    `json:"enable_forward,omitempty"`
	UpdateMulti              *bool    `json:"update_multi,omitempty"`
	WidthMode                string   `json:"width_mode,omitempty"`
	UseCustomTranslation     *bool    `json:"use_custom_translation,omitempty"`
	EnableForwardInteraction *bool    `json:"enable_forward_interaction,omitempty"`
	Style                    Style    `json:"style,omitempty"`
}

// Card represents the entire JSON structure of a card.
type Card struct {
	Config       *Config       `json:"config"`
	CardLink     *CardLink     `json:"card_link,omitempty"`
	I18nElements *I18nElements `json:"i18n_elements,omitempty"`
	Header       *Header       `json:"header,omitempty"`
}

// Template represents theme styles enumeration.
type Template string

const (
	ThemeBlue      Template = "blue"
	ThemeWathet    Template = "wathet"
	ThemeTurquoise Template = "turquoise"
	ThemeGreen     Template = "green"
	ThemeYellow    Template = "yellow"
	ThemeOrange    Template = "orange"
	ThemeRed       Template = "red"
	ThemeCarmine   Template = "carmine"
	ThemeViolet    Template = "violet"
	ThemePurple    Template = "purple"
	ThemeIndigo    Template = "indigo"
	ThemeGrey      Template = "grey"
	ThemeDefault   Template = "default"
)

// Color represents color effects enumeration.
type Color string

const (
	ColorNeutral   Color = "neutral"
	ColorBlue      Color = "blue"
	ColorTurquoise Color = "turquoise"
	ColorLime      Color = "lime"
	ColorOrange    Color = "orange"
	ColorViolet    Color = "violet"
	ColorIndigo    Color = "indigo"
	ColorWathet    Color = "wathet"
	ColorGreen     Color = "green"
	ColorYellow    Color = "yellow"
	ColorRed       Color = "red"
	ColorPurple    Color = "purple"
	ColorCarmine   Color = "carmine"
)
