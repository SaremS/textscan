#include "parser.h"

#if defined(__GNUC__) || defined(__clang__)
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wmissing-field-initializers"
#endif

#define LANGUAGE_VERSION 14
#define STATE_COUNT 14
#define LARGE_STATE_COUNT 5
#define SYMBOL_COUNT 8
#define ALIAS_COUNT 0
#define TOKEN_COUNT 4
#define EXTERNAL_TOKEN_COUNT 0
#define FIELD_COUNT 0
#define MAX_ALIAS_SEQUENCE_LENGTH 3
#define PRODUCTION_ID_COUNT 1

enum {
  sym_sentence = 1,
  sym__blank_line = 2,
  sym__line_break = 3,
  sym_source_file = 4,
  sym_paragraph = 5,
  aux_sym_source_file_repeat1 = 6,
  aux_sym_paragraph_repeat1 = 7,
};

static const char * const ts_symbol_names[] = {
  [ts_builtin_sym_end] = "end",
  [sym_sentence] = "sentence",
  [sym__blank_line] = "_blank_line",
  [sym__line_break] = "_line_break",
  [sym_source_file] = "source_file",
  [sym_paragraph] = "paragraph",
  [aux_sym_source_file_repeat1] = "source_file_repeat1",
  [aux_sym_paragraph_repeat1] = "paragraph_repeat1",
};

static const TSSymbol ts_symbol_map[] = {
  [ts_builtin_sym_end] = ts_builtin_sym_end,
  [sym_sentence] = sym_sentence,
  [sym__blank_line] = sym__blank_line,
  [sym__line_break] = sym__line_break,
  [sym_source_file] = sym_source_file,
  [sym_paragraph] = sym_paragraph,
  [aux_sym_source_file_repeat1] = aux_sym_source_file_repeat1,
  [aux_sym_paragraph_repeat1] = aux_sym_paragraph_repeat1,
};

static const TSSymbolMetadata ts_symbol_metadata[] = {
  [ts_builtin_sym_end] = {
    .visible = false,
    .named = true,
  },
  [sym_sentence] = {
    .visible = true,
    .named = true,
  },
  [sym__blank_line] = {
    .visible = false,
    .named = true,
  },
  [sym__line_break] = {
    .visible = false,
    .named = true,
  },
  [sym_source_file] = {
    .visible = true,
    .named = true,
  },
  [sym_paragraph] = {
    .visible = true,
    .named = true,
  },
  [aux_sym_source_file_repeat1] = {
    .visible = false,
    .named = false,
  },
  [aux_sym_paragraph_repeat1] = {
    .visible = false,
    .named = false,
  },
};

static const TSSymbol ts_alias_sequences[PRODUCTION_ID_COUNT][MAX_ALIAS_SEQUENCE_LENGTH] = {
  [0] = {0},
};

static const uint16_t ts_non_terminal_alias_map[] = {
  0,
};

static const TSStateId ts_primary_state_ids[STATE_COUNT] = {
  [0] = 0,
  [1] = 1,
  [2] = 2,
  [3] = 3,
  [4] = 4,
  [5] = 5,
  [6] = 6,
  [7] = 7,
  [8] = 8,
  [9] = 9,
  [10] = 10,
  [11] = 11,
  [12] = 12,
  [13] = 13,
};

static bool ts_lex(TSLexer *lexer, TSStateId state) {
  START_LEXER();
  eof = lexer->eof(lexer);
  switch (state) {
    case 0:
      if (eof) ADVANCE(9);
      if (lookahead == '\n') ADVANCE(16);
      if (lookahead == '\r') ADVANCE(1);
      if (lookahead != 0 &&
          lookahead != '!' &&
          lookahead != '.' &&
          lookahead != '?') ADVANCE(11);
      END_STATE();
    case 1:
      if (lookahead == '\n') ADVANCE(16);
      END_STATE();
    case 2:
      if (lookahead == '\n') ADVANCE(14);
      END_STATE();
    case 3:
      if (lookahead == '\n') ADVANCE(14);
      if (lookahead == '\r') ADVANCE(2);
      if (lookahead == '\t' ||
          lookahead == ' ') ADVANCE(3);
      END_STATE();
    case 4:
      if (lookahead == '\n') ADVANCE(7);
      END_STATE();
    case 5:
      if (lookahead == '\n') ADVANCE(15);
      END_STATE();
    case 6:
      if (lookahead == '\n') ADVANCE(17);
      END_STATE();
    case 7:
      if (lookahead != 0 &&
          lookahead != '\n' &&
          lookahead != '\r' &&
          lookahead != '!' &&
          lookahead != '.' &&
          lookahead != '?') ADVANCE(11);
      END_STATE();
    case 8:
      if (eof) ADVANCE(9);
      if (lookahead == '\n') ADVANCE(17);
      if (lookahead == '\r') ADVANCE(6);
      END_STATE();
    case 9:
      ACCEPT_TOKEN(ts_builtin_sym_end);
      END_STATE();
    case 10:
      ACCEPT_TOKEN(sym_sentence);
      if (lookahead == '\n') ADVANCE(7);
      if (lookahead == '\r') ADVANCE(4);
      if (lookahead == '\t' ||
          lookahead == ' ') ADVANCE(10);
      if (lookahead == '!' ||
          lookahead == '.' ||
          lookahead == '?') ADVANCE(13);
      if (lookahead != 0) ADVANCE(11);
      END_STATE();
    case 11:
      ACCEPT_TOKEN(sym_sentence);
      if (lookahead == '\n') ADVANCE(7);
      if (lookahead == '\r') ADVANCE(4);
      if (lookahead == '!' ||
          lookahead == '.' ||
          lookahead == '?') ADVANCE(13);
      if (lookahead != 0) ADVANCE(11);
      END_STATE();
    case 12:
      ACCEPT_TOKEN(sym_sentence);
      if (lookahead == '\n') ADVANCE(15);
      if (lookahead == '\r') ADVANCE(5);
      if (lookahead == '\t' ||
          lookahead == ' ') ADVANCE(12);
      if (lookahead == '!' ||
          lookahead == '.' ||
          lookahead == '?') ADVANCE(13);
      if (lookahead != 0) ADVANCE(11);
      END_STATE();
    case 13:
      ACCEPT_TOKEN(sym_sentence);
      if (lookahead == '!' ||
          lookahead == '.' ||
          lookahead == '?') ADVANCE(13);
      END_STATE();
    case 14:
      ACCEPT_TOKEN(sym__blank_line);
      if (lookahead == '\t' ||
          lookahead == ' ') ADVANCE(14);
      END_STATE();
    case 15:
      ACCEPT_TOKEN(sym__blank_line);
      if (lookahead == '\t' ||
          lookahead == ' ') ADVANCE(10);
      if (lookahead != 0 &&
          lookahead != '\n' &&
          lookahead != '\r' &&
          lookahead != '!' &&
          lookahead != '.' &&
          lookahead != '?') ADVANCE(11);
      END_STATE();
    case 16:
      ACCEPT_TOKEN(sym__line_break);
      if (lookahead == '\n') ADVANCE(14);
      if (lookahead == '\r') ADVANCE(2);
      if (lookahead == '\t' ||
          lookahead == ' ') ADVANCE(12);
      if (lookahead != 0 &&
          lookahead != '!' &&
          lookahead != '.' &&
          lookahead != '?') ADVANCE(11);
      END_STATE();
    case 17:
      ACCEPT_TOKEN(sym__line_break);
      if (lookahead == '\n') ADVANCE(14);
      if (lookahead == '\r') ADVANCE(2);
      if (lookahead == '\t' ||
          lookahead == ' ') ADVANCE(3);
      END_STATE();
    default:
      return false;
  }
}

static const TSLexMode ts_lex_modes[STATE_COUNT] = {
  [0] = {.lex_state = 0},
  [1] = {.lex_state = 0},
  [2] = {.lex_state = 0},
  [3] = {.lex_state = 0},
  [4] = {.lex_state = 0},
  [5] = {.lex_state = 8},
  [6] = {.lex_state = 8},
  [7] = {.lex_state = 8},
  [8] = {.lex_state = 8},
  [9] = {.lex_state = 8},
  [10] = {.lex_state = 0},
  [11] = {.lex_state = 0},
  [12] = {.lex_state = 0},
  [13] = {.lex_state = 0},
};

static const uint16_t ts_parse_table[LARGE_STATE_COUNT][SYMBOL_COUNT] = {
  [0] = {
    [ts_builtin_sym_end] = ACTIONS(1),
    [sym_sentence] = ACTIONS(1),
    [sym__blank_line] = ACTIONS(1),
    [sym__line_break] = ACTIONS(1),
  },
  [1] = {
    [sym_source_file] = STATE(11),
    [sym_paragraph] = STATE(5),
    [aux_sym_source_file_repeat1] = STATE(6),
    [aux_sym_paragraph_repeat1] = STATE(3),
    [ts_builtin_sym_end] = ACTIONS(3),
    [sym_sentence] = ACTIONS(5),
    [sym__blank_line] = ACTIONS(7),
    [sym__line_break] = ACTIONS(9),
  },
  [2] = {
    [sym_paragraph] = STATE(9),
    [aux_sym_paragraph_repeat1] = STATE(3),
    [ts_builtin_sym_end] = ACTIONS(11),
    [sym_sentence] = ACTIONS(5),
    [sym__blank_line] = ACTIONS(13),
    [sym__line_break] = ACTIONS(13),
  },
  [3] = {
    [aux_sym_paragraph_repeat1] = STATE(4),
    [ts_builtin_sym_end] = ACTIONS(15),
    [sym_sentence] = ACTIONS(17),
    [sym__blank_line] = ACTIONS(19),
    [sym__line_break] = ACTIONS(19),
  },
  [4] = {
    [aux_sym_paragraph_repeat1] = STATE(4),
    [ts_builtin_sym_end] = ACTIONS(21),
    [sym_sentence] = ACTIONS(23),
    [sym__blank_line] = ACTIONS(26),
    [sym__line_break] = ACTIONS(26),
  },
};

static const uint16_t ts_small_parse_table[] = {
  [0] = 4,
    ACTIONS(28), 1,
      ts_builtin_sym_end,
    ACTIONS(30), 1,
      sym__blank_line,
    ACTIONS(32), 1,
      sym__line_break,
    STATE(7), 1,
      aux_sym_source_file_repeat1,
  [13] = 4,
    ACTIONS(28), 1,
      ts_builtin_sym_end,
    ACTIONS(30), 1,
      sym__blank_line,
    ACTIONS(32), 1,
      sym__line_break,
    STATE(8), 1,
      aux_sym_source_file_repeat1,
  [26] = 4,
    ACTIONS(30), 1,
      sym__blank_line,
    ACTIONS(34), 1,
      ts_builtin_sym_end,
    ACTIONS(36), 1,
      sym__line_break,
    STATE(8), 1,
      aux_sym_source_file_repeat1,
  [39] = 4,
    ACTIONS(38), 1,
      ts_builtin_sym_end,
    ACTIONS(40), 1,
      sym__blank_line,
    ACTIONS(43), 1,
      sym__line_break,
    STATE(8), 1,
      aux_sym_source_file_repeat1,
  [52] = 2,
    ACTIONS(43), 1,
      sym__line_break,
    ACTIONS(38), 2,
      ts_builtin_sym_end,
      sym__blank_line,
  [60] = 1,
    ACTIONS(28), 1,
      ts_builtin_sym_end,
  [64] = 1,
    ACTIONS(45), 1,
      ts_builtin_sym_end,
  [68] = 1,
    ACTIONS(34), 1,
      ts_builtin_sym_end,
  [72] = 1,
    ACTIONS(47), 1,
      ts_builtin_sym_end,
};

static const uint32_t ts_small_parse_table_map[] = {
  [SMALL_STATE(5)] = 0,
  [SMALL_STATE(6)] = 13,
  [SMALL_STATE(7)] = 26,
  [SMALL_STATE(8)] = 39,
  [SMALL_STATE(9)] = 52,
  [SMALL_STATE(10)] = 60,
  [SMALL_STATE(11)] = 64,
  [SMALL_STATE(12)] = 68,
  [SMALL_STATE(13)] = 72,
};

static const TSParseActionEntry ts_parse_actions[] = {
  [0] = {.entry = {.count = 0, .reusable = false}},
  [1] = {.entry = {.count = 1, .reusable = false}}, RECOVER(),
  [3] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_source_file, 0),
  [5] = {.entry = {.count = 1, .reusable = true}}, SHIFT(3),
  [7] = {.entry = {.count = 1, .reusable = false}}, SHIFT(2),
  [9] = {.entry = {.count = 1, .reusable = false}}, SHIFT(10),
  [11] = {.entry = {.count = 1, .reusable = true}}, REDUCE(aux_sym_source_file_repeat1, 1),
  [13] = {.entry = {.count = 1, .reusable = false}}, REDUCE(aux_sym_source_file_repeat1, 1),
  [15] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_paragraph, 1),
  [17] = {.entry = {.count = 1, .reusable = true}}, SHIFT(4),
  [19] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_paragraph, 1),
  [21] = {.entry = {.count = 1, .reusable = true}}, REDUCE(aux_sym_paragraph_repeat1, 2),
  [23] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_paragraph_repeat1, 2), SHIFT_REPEAT(4),
  [26] = {.entry = {.count = 1, .reusable = false}}, REDUCE(aux_sym_paragraph_repeat1, 2),
  [28] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_source_file, 1),
  [30] = {.entry = {.count = 1, .reusable = true}}, SHIFT(2),
  [32] = {.entry = {.count = 1, .reusable = false}}, SHIFT(12),
  [34] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_source_file, 2),
  [36] = {.entry = {.count = 1, .reusable = false}}, SHIFT(13),
  [38] = {.entry = {.count = 1, .reusable = true}}, REDUCE(aux_sym_source_file_repeat1, 2),
  [40] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_source_file_repeat1, 2), SHIFT_REPEAT(2),
  [43] = {.entry = {.count = 1, .reusable = false}}, REDUCE(aux_sym_source_file_repeat1, 2),
  [45] = {.entry = {.count = 1, .reusable = true}},  ACCEPT_INPUT(),
  [47] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_source_file, 3),
};

#ifdef __cplusplus
extern "C" {
#endif
#ifdef _WIN32
#define extern __declspec(dllexport)
#endif

extern const TSLanguage *tree_sitter_plaintext(void) {
  static const TSLanguage language = {
    .version = LANGUAGE_VERSION,
    .symbol_count = SYMBOL_COUNT,
    .alias_count = ALIAS_COUNT,
    .token_count = TOKEN_COUNT,
    .external_token_count = EXTERNAL_TOKEN_COUNT,
    .state_count = STATE_COUNT,
    .large_state_count = LARGE_STATE_COUNT,
    .production_id_count = PRODUCTION_ID_COUNT,
    .field_count = FIELD_COUNT,
    .max_alias_sequence_length = MAX_ALIAS_SEQUENCE_LENGTH,
    .parse_table = &ts_parse_table[0][0],
    .small_parse_table = ts_small_parse_table,
    .small_parse_table_map = ts_small_parse_table_map,
    .parse_actions = ts_parse_actions,
    .symbol_names = ts_symbol_names,
    .symbol_metadata = ts_symbol_metadata,
    .public_symbol_map = ts_symbol_map,
    .alias_map = ts_non_terminal_alias_map,
    .alias_sequences = &ts_alias_sequences[0][0],
    .lex_modes = ts_lex_modes,
    .lex_fn = ts_lex,
    .primary_state_ids = ts_primary_state_ids,
  };
  return &language;
}
#ifdef __cplusplus
}
#endif
