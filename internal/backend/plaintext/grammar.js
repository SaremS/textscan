module.exports = grammar({
  name: 'plaintext',

  extras: () => [],

  rules: {
    source_file: $ => seq(
      optional($.paragraph),
      repeat(seq($._blank_line, optional($.paragraph))),
      optional($._line_break),
    ),

    paragraph: $ => repeat1($.sentence),

    sentence: () => token(choice(
      /([^.!?\r\n]|\r?\n[^\r\n.!?])+[.!?]+/,
      /([^.!?\r\n]|\r?\n[^\r\n.!?])+/,
    )),

    _blank_line: () => token(/(\r?\n[ \t]*){2,}/),
    _line_break: () => token(/\r?\n/),
  },
});
