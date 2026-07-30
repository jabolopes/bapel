
// Generated from cpp_parser/bapel.g4 by ANTLR 4.9.2


#include "bapelVisitor.h"

#include "bapelParser.h"


using namespace antlrcpp;
using namespace antlr4;

bapelParser::bapelParser(TokenStream *input) : Parser(input) {
  _interpreter = new atn::ParserATNSimulator(this, _atn, _decisionToDFA, _sharedContextCache);
}

bapelParser::~bapelParser() {
  delete _interpreter;
}

std::string bapelParser::getGrammarFileName() const {
  return "bapel.g4";
}

const std::vector<std::string>& bapelParser::getRuleNames() const {
  return _ruleNames;
}

dfa::Vocabulary& bapelParser::getVocabulary() const {
  return _vocabulary;
}


//----------------- SourceFileContext ------------------------------------------------------------------

bapelParser::SourceFileContext::SourceFileContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}


size_t bapelParser::SourceFileContext::getRuleIndex() const {
  return bapelParser::RuleSourceFile;
}

void bapelParser::SourceFileContext::copyFrom(SourceFileContext *ctx) {
  ParserRuleContext::copyFrom(ctx);
}

//----------------- ImplSourceFileContext ------------------------------------------------------------------

bapelParser::ImplementsHeaderContext* bapelParser::ImplSourceFileContext::implementsHeader() {
  return getRuleContext<bapelParser::ImplementsHeaderContext>(0);
}

tree::TerminalNode* bapelParser::ImplSourceFileContext::EOF() {
  return getToken(bapelParser::EOF, 0);
}

bapelParser::ImportsSectionContext* bapelParser::ImplSourceFileContext::importsSection() {
  return getRuleContext<bapelParser::ImportsSectionContext>(0);
}

bapelParser::ImplsSectionContext* bapelParser::ImplSourceFileContext::implsSection() {
  return getRuleContext<bapelParser::ImplsSectionContext>(0);
}

bapelParser::FlagsSectionContext* bapelParser::ImplSourceFileContext::flagsSection() {
  return getRuleContext<bapelParser::FlagsSectionContext>(0);
}

bapelParser::SourcesContext* bapelParser::ImplSourceFileContext::sources() {
  return getRuleContext<bapelParser::SourcesContext>(0);
}

bapelParser::ImplSourceFileContext::ImplSourceFileContext(SourceFileContext *ctx) { copyFrom(ctx); }


antlrcpp::Any bapelParser::ImplSourceFileContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitImplSourceFile(this);
  else
    return visitor->visitChildren(this);
}
//----------------- BaseSourceFileContext ------------------------------------------------------------------

bapelParser::ModuleHeaderContext* bapelParser::BaseSourceFileContext::moduleHeader() {
  return getRuleContext<bapelParser::ModuleHeaderContext>(0);
}

tree::TerminalNode* bapelParser::BaseSourceFileContext::EOF() {
  return getToken(bapelParser::EOF, 0);
}

bapelParser::ImportsSectionContext* bapelParser::BaseSourceFileContext::importsSection() {
  return getRuleContext<bapelParser::ImportsSectionContext>(0);
}

bapelParser::ImplsSectionContext* bapelParser::BaseSourceFileContext::implsSection() {
  return getRuleContext<bapelParser::ImplsSectionContext>(0);
}

bapelParser::FlagsSectionContext* bapelParser::BaseSourceFileContext::flagsSection() {
  return getRuleContext<bapelParser::FlagsSectionContext>(0);
}

bapelParser::SourcesContext* bapelParser::BaseSourceFileContext::sources() {
  return getRuleContext<bapelParser::SourcesContext>(0);
}

bapelParser::BaseSourceFileContext::BaseSourceFileContext(SourceFileContext *ctx) { copyFrom(ctx); }


antlrcpp::Any bapelParser::BaseSourceFileContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitBaseSourceFile(this);
  else
    return visitor->visitChildren(this);
}
bapelParser::SourceFileContext* bapelParser::sourceFile() {
  SourceFileContext *_localctx = _tracker.createInstance<SourceFileContext>(_ctx, getState());
  enterRule(_localctx, 0, bapelParser::RuleSourceFile);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(202);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::MODULE: {
        _localctx = dynamic_cast<SourceFileContext *>(_tracker.createInstance<bapelParser::BaseSourceFileContext>(_localctx));
        enterOuterAlt(_localctx, 1);
        setState(172);
        moduleHeader();
        setState(174);
        _errHandler->sync(this);

        _la = _input->LA(1);
        if (_la == bapelParser::IMPORTS) {
          setState(173);
          importsSection();
        }
        setState(177);
        _errHandler->sync(this);

        _la = _input->LA(1);
        if (_la == bapelParser::IMPLS) {
          setState(176);
          implsSection();
        }
        setState(180);
        _errHandler->sync(this);

        _la = _input->LA(1);
        if (_la == bapelParser::FLAGS) {
          setState(179);
          flagsSection();
        }
        setState(183);
        _errHandler->sync(this);

        _la = _input->LA(1);
        if ((((_la & ~ 0x3fULL) == 0) &&
          ((1ULL << _la) & ((1ULL << bapelParser::PUB)
          | (1ULL << bapelParser::DECL)
          | (1ULL << bapelParser::FN)
          | (1ULL << bapelParser::TYPE)
          | (1ULL << bapelParser::TRAIT)
          | (1ULL << bapelParser::IMPL))) != 0)) {
          setState(182);
          sources();
        }
        setState(185);
        match(bapelParser::EOF);
        break;
      }

      case bapelParser::IMPLEMENTS: {
        _localctx = dynamic_cast<SourceFileContext *>(_tracker.createInstance<bapelParser::ImplSourceFileContext>(_localctx));
        enterOuterAlt(_localctx, 2);
        setState(187);
        implementsHeader();
        setState(189);
        _errHandler->sync(this);

        _la = _input->LA(1);
        if (_la == bapelParser::IMPORTS) {
          setState(188);
          importsSection();
        }
        setState(192);
        _errHandler->sync(this);

        _la = _input->LA(1);
        if (_la == bapelParser::IMPLS) {
          setState(191);
          implsSection();
        }
        setState(195);
        _errHandler->sync(this);

        _la = _input->LA(1);
        if (_la == bapelParser::FLAGS) {
          setState(194);
          flagsSection();
        }
        setState(198);
        _errHandler->sync(this);

        _la = _input->LA(1);
        if ((((_la & ~ 0x3fULL) == 0) &&
          ((1ULL << _la) & ((1ULL << bapelParser::PUB)
          | (1ULL << bapelParser::DECL)
          | (1ULL << bapelParser::FN)
          | (1ULL << bapelParser::TYPE)
          | (1ULL << bapelParser::TRAIT)
          | (1ULL << bapelParser::IMPL))) != 0)) {
          setState(197);
          sources();
        }
        setState(200);
        match(bapelParser::EOF);
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ModuleHeaderContext ------------------------------------------------------------------

bapelParser::ModuleHeaderContext::ModuleHeaderContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::ModuleHeaderContext::MODULE() {
  return getToken(bapelParser::MODULE, 0);
}

bapelParser::ModuleIDContext* bapelParser::ModuleHeaderContext::moduleID() {
  return getRuleContext<bapelParser::ModuleIDContext>(0);
}


size_t bapelParser::ModuleHeaderContext::getRuleIndex() const {
  return bapelParser::RuleModuleHeader;
}


antlrcpp::Any bapelParser::ModuleHeaderContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitModuleHeader(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::ModuleHeaderContext* bapelParser::moduleHeader() {
  ModuleHeaderContext *_localctx = _tracker.createInstance<ModuleHeaderContext>(_ctx, getState());
  enterRule(_localctx, 2, bapelParser::RuleModuleHeader);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(204);
    match(bapelParser::MODULE);
    setState(205);
    moduleID();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ImplementsHeaderContext ------------------------------------------------------------------

bapelParser::ImplementsHeaderContext::ImplementsHeaderContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::ImplementsHeaderContext::IMPLEMENTS() {
  return getToken(bapelParser::IMPLEMENTS, 0);
}

bapelParser::ModuleIDContext* bapelParser::ImplementsHeaderContext::moduleID() {
  return getRuleContext<bapelParser::ModuleIDContext>(0);
}


size_t bapelParser::ImplementsHeaderContext::getRuleIndex() const {
  return bapelParser::RuleImplementsHeader;
}


antlrcpp::Any bapelParser::ImplementsHeaderContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitImplementsHeader(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::ImplementsHeaderContext* bapelParser::implementsHeader() {
  ImplementsHeaderContext *_localctx = _tracker.createInstance<ImplementsHeaderContext>(_ctx, getState());
  enterRule(_localctx, 4, bapelParser::RuleImplementsHeader);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(207);
    match(bapelParser::IMPLEMENTS);
    setState(208);
    moduleID();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- WorkspaceContext ------------------------------------------------------------------

bapelParser::WorkspaceContext::WorkspaceContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::WorkspaceContext::WORKSPACE() {
  return getToken(bapelParser::WORKSPACE, 0);
}

tree::TerminalNode* bapelParser::WorkspaceContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

bapelParser::PackagesSectionContext* bapelParser::WorkspaceContext::packagesSection() {
  return getRuleContext<bapelParser::PackagesSectionContext>(0);
}

tree::TerminalNode* bapelParser::WorkspaceContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}


size_t bapelParser::WorkspaceContext::getRuleIndex() const {
  return bapelParser::RuleWorkspace;
}


antlrcpp::Any bapelParser::WorkspaceContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitWorkspace(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::WorkspaceContext* bapelParser::workspace() {
  WorkspaceContext *_localctx = _tracker.createInstance<WorkspaceContext>(_ctx, getState());
  enterRule(_localctx, 6, bapelParser::RuleWorkspace);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(210);
    match(bapelParser::WORKSPACE);
    setState(211);
    match(bapelParser::LBRACE);
    setState(212);
    packagesSection();
    setState(213);
    match(bapelParser::RBRACE);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- PackagesSectionContext ------------------------------------------------------------------

bapelParser::PackagesSectionContext::PackagesSectionContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::PackagesSectionContext::PACKAGES() {
  return getToken(bapelParser::PACKAGES, 0);
}

tree::TerminalNode* bapelParser::PackagesSectionContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

tree::TerminalNode* bapelParser::PackagesSectionContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}

std::vector<bapelParser::PackageRuleContext *> bapelParser::PackagesSectionContext::packageRule() {
  return getRuleContexts<bapelParser::PackageRuleContext>();
}

bapelParser::PackageRuleContext* bapelParser::PackagesSectionContext::packageRule(size_t i) {
  return getRuleContext<bapelParser::PackageRuleContext>(i);
}


size_t bapelParser::PackagesSectionContext::getRuleIndex() const {
  return bapelParser::RulePackagesSection;
}


antlrcpp::Any bapelParser::PackagesSectionContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitPackagesSection(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::PackagesSectionContext* bapelParser::packagesSection() {
  PackagesSectionContext *_localctx = _tracker.createInstance<PackagesSectionContext>(_ctx, getState());
  enterRule(_localctx, 8, bapelParser::RulePackagesSection);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(215);
    match(bapelParser::PACKAGES);
    setState(216);
    match(bapelParser::LBRACE);
    setState(218); 
    _errHandler->sync(this);
    _la = _input->LA(1);
    do {
      setState(217);
      packageRule();
      setState(220); 
      _errHandler->sync(this);
      _la = _input->LA(1);
    } while (_la == bapelParser::PREFIX

    || _la == bapelParser::MODULE);
    setState(222);
    match(bapelParser::RBRACE);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- PackageRuleContext ------------------------------------------------------------------

bapelParser::PackageRuleContext::PackageRuleContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::PackageRuleContext::PREFIX() {
  return getToken(bapelParser::PREFIX, 0);
}

bapelParser::ModuleIDContext* bapelParser::PackageRuleContext::moduleID() {
  return getRuleContext<bapelParser::ModuleIDContext>(0);
}

tree::TerminalNode* bapelParser::PackageRuleContext::IN() {
  return getToken(bapelParser::IN, 0);
}

bapelParser::FilenameContext* bapelParser::PackageRuleContext::filename() {
  return getRuleContext<bapelParser::FilenameContext>(0);
}

tree::TerminalNode* bapelParser::PackageRuleContext::MODULE() {
  return getToken(bapelParser::MODULE, 0);
}


size_t bapelParser::PackageRuleContext::getRuleIndex() const {
  return bapelParser::RulePackageRule;
}


antlrcpp::Any bapelParser::PackageRuleContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitPackageRule(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::PackageRuleContext* bapelParser::packageRule() {
  PackageRuleContext *_localctx = _tracker.createInstance<PackageRuleContext>(_ctx, getState());
  enterRule(_localctx, 10, bapelParser::RulePackageRule);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(234);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::PREFIX: {
        enterOuterAlt(_localctx, 1);
        setState(224);
        match(bapelParser::PREFIX);
        setState(225);
        moduleID();
        setState(226);
        match(bapelParser::IN);
        setState(227);
        filename();
        break;
      }

      case bapelParser::MODULE: {
        enterOuterAlt(_localctx, 2);
        setState(229);
        match(bapelParser::MODULE);
        setState(230);
        moduleID();
        setState(231);
        match(bapelParser::IN);
        setState(232);
        filename();
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ImportsSectionContext ------------------------------------------------------------------

bapelParser::ImportsSectionContext::ImportsSectionContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::ImportsSectionContext::IMPORTS() {
  return getToken(bapelParser::IMPORTS, 0);
}

tree::TerminalNode* bapelParser::ImportsSectionContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

tree::TerminalNode* bapelParser::ImportsSectionContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}

std::vector<bapelParser::ModuleIDContext *> bapelParser::ImportsSectionContext::moduleID() {
  return getRuleContexts<bapelParser::ModuleIDContext>();
}

bapelParser::ModuleIDContext* bapelParser::ImportsSectionContext::moduleID(size_t i) {
  return getRuleContext<bapelParser::ModuleIDContext>(i);
}


size_t bapelParser::ImportsSectionContext::getRuleIndex() const {
  return bapelParser::RuleImportsSection;
}


antlrcpp::Any bapelParser::ImportsSectionContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitImportsSection(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::ImportsSectionContext* bapelParser::importsSection() {
  ImportsSectionContext *_localctx = _tracker.createInstance<ImportsSectionContext>(_ctx, getState());
  enterRule(_localctx, 12, bapelParser::RuleImportsSection);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(236);
    match(bapelParser::IMPORTS);
    setState(237);
    match(bapelParser::LBRACE);
    setState(239); 
    _errHandler->sync(this);
    _la = _input->LA(1);
    do {
      setState(238);
      moduleID();
      setState(241); 
      _errHandler->sync(this);
      _la = _input->LA(1);
    } while (_la == bapelParser::IDENTIFIER);
    setState(243);
    match(bapelParser::RBRACE);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ImplsSectionContext ------------------------------------------------------------------

bapelParser::ImplsSectionContext::ImplsSectionContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::ImplsSectionContext::IMPLS() {
  return getToken(bapelParser::IMPLS, 0);
}

tree::TerminalNode* bapelParser::ImplsSectionContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

tree::TerminalNode* bapelParser::ImplsSectionContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}

std::vector<bapelParser::FilenameContext *> bapelParser::ImplsSectionContext::filename() {
  return getRuleContexts<bapelParser::FilenameContext>();
}

bapelParser::FilenameContext* bapelParser::ImplsSectionContext::filename(size_t i) {
  return getRuleContext<bapelParser::FilenameContext>(i);
}


size_t bapelParser::ImplsSectionContext::getRuleIndex() const {
  return bapelParser::RuleImplsSection;
}


antlrcpp::Any bapelParser::ImplsSectionContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitImplsSection(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::ImplsSectionContext* bapelParser::implsSection() {
  ImplsSectionContext *_localctx = _tracker.createInstance<ImplsSectionContext>(_ctx, getState());
  enterRule(_localctx, 14, bapelParser::RuleImplsSection);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(245);
    match(bapelParser::IMPLS);
    setState(246);
    match(bapelParser::LBRACE);
    setState(248); 
    _errHandler->sync(this);
    _la = _input->LA(1);
    do {
      setState(247);
      filename();
      setState(250); 
      _errHandler->sync(this);
      _la = _input->LA(1);
    } while (_la == bapelParser::STRING_LITERAL);
    setState(252);
    match(bapelParser::RBRACE);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- FlagsSectionContext ------------------------------------------------------------------

bapelParser::FlagsSectionContext::FlagsSectionContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::FlagsSectionContext::FLAGS() {
  return getToken(bapelParser::FLAGS, 0);
}

tree::TerminalNode* bapelParser::FlagsSectionContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

tree::TerminalNode* bapelParser::FlagsSectionContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}

std::vector<bapelParser::FilenameContext *> bapelParser::FlagsSectionContext::filename() {
  return getRuleContexts<bapelParser::FilenameContext>();
}

bapelParser::FilenameContext* bapelParser::FlagsSectionContext::filename(size_t i) {
  return getRuleContext<bapelParser::FilenameContext>(i);
}


size_t bapelParser::FlagsSectionContext::getRuleIndex() const {
  return bapelParser::RuleFlagsSection;
}


antlrcpp::Any bapelParser::FlagsSectionContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitFlagsSection(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::FlagsSectionContext* bapelParser::flagsSection() {
  FlagsSectionContext *_localctx = _tracker.createInstance<FlagsSectionContext>(_ctx, getState());
  enterRule(_localctx, 16, bapelParser::RuleFlagsSection);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(254);
    match(bapelParser::FLAGS);
    setState(255);
    match(bapelParser::LBRACE);
    setState(257); 
    _errHandler->sync(this);
    _la = _input->LA(1);
    do {
      setState(256);
      filename();
      setState(259); 
      _errHandler->sync(this);
      _la = _input->LA(1);
    } while (_la == bapelParser::STRING_LITERAL);
    setState(261);
    match(bapelParser::RBRACE);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ModuleIDContext ------------------------------------------------------------------

bapelParser::ModuleIDContext::ModuleIDContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

std::vector<tree::TerminalNode *> bapelParser::ModuleIDContext::IDENTIFIER() {
  return getTokens(bapelParser::IDENTIFIER);
}

tree::TerminalNode* bapelParser::ModuleIDContext::IDENTIFIER(size_t i) {
  return getToken(bapelParser::IDENTIFIER, i);
}

std::vector<tree::TerminalNode *> bapelParser::ModuleIDContext::DOT() {
  return getTokens(bapelParser::DOT);
}

tree::TerminalNode* bapelParser::ModuleIDContext::DOT(size_t i) {
  return getToken(bapelParser::DOT, i);
}


size_t bapelParser::ModuleIDContext::getRuleIndex() const {
  return bapelParser::RuleModuleID;
}


antlrcpp::Any bapelParser::ModuleIDContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitModuleID(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::ModuleIDContext* bapelParser::moduleID() {
  ModuleIDContext *_localctx = _tracker.createInstance<ModuleIDContext>(_ctx, getState());
  enterRule(_localctx, 18, bapelParser::RuleModuleID);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(263);
    match(bapelParser::IDENTIFIER);
    setState(268);
    _errHandler->sync(this);
    _la = _input->LA(1);
    while (_la == bapelParser::DOT) {
      setState(264);
      match(bapelParser::DOT);
      setState(265);
      match(bapelParser::IDENTIFIER);
      setState(270);
      _errHandler->sync(this);
      _la = _input->LA(1);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- FilenameContext ------------------------------------------------------------------

bapelParser::FilenameContext::FilenameContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::FilenameContext::STRING_LITERAL() {
  return getToken(bapelParser::STRING_LITERAL, 0);
}


size_t bapelParser::FilenameContext::getRuleIndex() const {
  return bapelParser::RuleFilename;
}


antlrcpp::Any bapelParser::FilenameContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitFilename(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::FilenameContext* bapelParser::filename() {
  FilenameContext *_localctx = _tracker.createInstance<FilenameContext>(_ctx, getState());
  enterRule(_localctx, 20, bapelParser::RuleFilename);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(271);
    match(bapelParser::STRING_LITERAL);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- SourcesContext ------------------------------------------------------------------

bapelParser::SourcesContext::SourcesContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

std::vector<bapelParser::SourceContext *> bapelParser::SourcesContext::source() {
  return getRuleContexts<bapelParser::SourceContext>();
}

bapelParser::SourceContext* bapelParser::SourcesContext::source(size_t i) {
  return getRuleContext<bapelParser::SourceContext>(i);
}


size_t bapelParser::SourcesContext::getRuleIndex() const {
  return bapelParser::RuleSources;
}


antlrcpp::Any bapelParser::SourcesContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitSources(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::SourcesContext* bapelParser::sources() {
  SourcesContext *_localctx = _tracker.createInstance<SourcesContext>(_ctx, getState());
  enterRule(_localctx, 22, bapelParser::RuleSources);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(274); 
    _errHandler->sync(this);
    _la = _input->LA(1);
    do {
      setState(273);
      source();
      setState(276); 
      _errHandler->sync(this);
      _la = _input->LA(1);
    } while ((((_la & ~ 0x3fULL) == 0) &&
      ((1ULL << _la) & ((1ULL << bapelParser::PUB)
      | (1ULL << bapelParser::DECL)
      | (1ULL << bapelParser::FN)
      | (1ULL << bapelParser::TYPE)
      | (1ULL << bapelParser::TRAIT)
      | (1ULL << bapelParser::IMPL))) != 0));
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- SourceContext ------------------------------------------------------------------

bapelParser::SourceContext::SourceContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::DeclNoExportContext* bapelParser::SourceContext::declNoExport() {
  return getRuleContext<bapelParser::DeclNoExportContext>(0);
}

bapelParser::FunctionNoExportContext* bapelParser::SourceContext::functionNoExport() {
  return getRuleContext<bapelParser::FunctionNoExportContext>(0);
}

tree::TerminalNode* bapelParser::SourceContext::PUB() {
  return getToken(bapelParser::PUB, 0);
}

bapelParser::TraitDeclContext* bapelParser::SourceContext::traitDecl() {
  return getRuleContext<bapelParser::TraitDeclContext>(0);
}

bapelParser::ImplBlockContext* bapelParser::SourceContext::implBlock() {
  return getRuleContext<bapelParser::ImplBlockContext>(0);
}


size_t bapelParser::SourceContext::getRuleIndex() const {
  return bapelParser::RuleSource;
}


antlrcpp::Any bapelParser::SourceContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitSource(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::SourceContext* bapelParser::source() {
  SourceContext *_localctx = _tracker.createInstance<SourceContext>(_ctx, getState());
  enterRule(_localctx, 24, bapelParser::RuleSource);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(286);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 16, _ctx)) {
    case 1: {
      enterOuterAlt(_localctx, 1);
      setState(278);
      declNoExport();
      break;
    }

    case 2: {
      enterOuterAlt(_localctx, 2);
      setState(279);
      functionNoExport();
      break;
    }

    case 3: {
      enterOuterAlt(_localctx, 3);
      setState(280);
      match(bapelParser::PUB);
      setState(281);
      functionNoExport();
      break;
    }

    case 4: {
      enterOuterAlt(_localctx, 4);
      setState(282);
      traitDecl();
      break;
    }

    case 5: {
      enterOuterAlt(_localctx, 5);
      setState(283);
      match(bapelParser::PUB);
      setState(284);
      traitDecl();
      break;
    }

    case 6: {
      enterOuterAlt(_localctx, 6);
      setState(285);
      implBlock();
      break;
    }

    default:
      break;
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- TraitDeclContext ------------------------------------------------------------------

bapelParser::TraitDeclContext::TraitDeclContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::TraitDeclContext::TRAIT() {
  return getToken(bapelParser::TRAIT, 0);
}

bapelParser::IdContext* bapelParser::TraitDeclContext::id() {
  return getRuleContext<bapelParser::IdContext>(0);
}

tree::TerminalNode* bapelParser::TraitDeclContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

tree::TerminalNode* bapelParser::TraitDeclContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}

bapelParser::TypeAbstractionContext* bapelParser::TraitDeclContext::typeAbstraction() {
  return getRuleContext<bapelParser::TypeAbstractionContext>(0);
}

std::vector<bapelParser::TraitMethodContext *> bapelParser::TraitDeclContext::traitMethod() {
  return getRuleContexts<bapelParser::TraitMethodContext>();
}

bapelParser::TraitMethodContext* bapelParser::TraitDeclContext::traitMethod(size_t i) {
  return getRuleContext<bapelParser::TraitMethodContext>(i);
}


size_t bapelParser::TraitDeclContext::getRuleIndex() const {
  return bapelParser::RuleTraitDecl;
}


antlrcpp::Any bapelParser::TraitDeclContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTraitDecl(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TraitDeclContext* bapelParser::traitDecl() {
  TraitDeclContext *_localctx = _tracker.createInstance<TraitDeclContext>(_ctx, getState());
  enterRule(_localctx, 26, bapelParser::RuleTraitDecl);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(288);
    match(bapelParser::TRAIT);
    setState(289);
    id();
    setState(291);
    _errHandler->sync(this);

    _la = _input->LA(1);
    if (_la == bapelParser::LBRACKET) {
      setState(290);
      typeAbstraction();
    }
    setState(293);
    match(bapelParser::LBRACE);
    setState(297);
    _errHandler->sync(this);
    _la = _input->LA(1);
    while (_la == bapelParser::FN) {
      setState(294);
      traitMethod();
      setState(299);
      _errHandler->sync(this);
      _la = _input->LA(1);
    }
    setState(300);
    match(bapelParser::RBRACE);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- TraitMethodContext ------------------------------------------------------------------

bapelParser::TraitMethodContext::TraitMethodContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::TraitMethodContext::FN() {
  return getToken(bapelParser::FN, 0);
}

bapelParser::IdContext* bapelParser::TraitMethodContext::id() {
  return getRuleContext<bapelParser::IdContext>(0);
}

bapelParser::FunctionArgsContext* bapelParser::TraitMethodContext::functionArgs() {
  return getRuleContext<bapelParser::FunctionArgsContext>(0);
}

tree::TerminalNode* bapelParser::TraitMethodContext::ARROW() {
  return getToken(bapelParser::ARROW, 0);
}

bapelParser::Type_Context* bapelParser::TraitMethodContext::type_() {
  return getRuleContext<bapelParser::Type_Context>(0);
}


size_t bapelParser::TraitMethodContext::getRuleIndex() const {
  return bapelParser::RuleTraitMethod;
}


antlrcpp::Any bapelParser::TraitMethodContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTraitMethod(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TraitMethodContext* bapelParser::traitMethod() {
  TraitMethodContext *_localctx = _tracker.createInstance<TraitMethodContext>(_ctx, getState());
  enterRule(_localctx, 28, bapelParser::RuleTraitMethod);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(302);
    match(bapelParser::FN);
    setState(303);
    id();
    setState(304);
    functionArgs();
    setState(305);
    match(bapelParser::ARROW);
    setState(306);
    type_();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ImplBlockContext ------------------------------------------------------------------

bapelParser::ImplBlockContext::ImplBlockContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}


size_t bapelParser::ImplBlockContext::getRuleIndex() const {
  return bapelParser::RuleImplBlock;
}

void bapelParser::ImplBlockContext::copyFrom(ImplBlockContext *ctx) {
  ParserRuleContext::copyFrom(ctx);
}

//----------------- InherentImplContext ------------------------------------------------------------------

tree::TerminalNode* bapelParser::InherentImplContext::IMPL() {
  return getToken(bapelParser::IMPL, 0);
}

bapelParser::Type_Context* bapelParser::InherentImplContext::type_() {
  return getRuleContext<bapelParser::Type_Context>(0);
}

tree::TerminalNode* bapelParser::InherentImplContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

tree::TerminalNode* bapelParser::InherentImplContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}

bapelParser::TypeAbstractionContext* bapelParser::InherentImplContext::typeAbstraction() {
  return getRuleContext<bapelParser::TypeAbstractionContext>(0);
}

std::vector<bapelParser::FunctionNoExportContext *> bapelParser::InherentImplContext::functionNoExport() {
  return getRuleContexts<bapelParser::FunctionNoExportContext>();
}

bapelParser::FunctionNoExportContext* bapelParser::InherentImplContext::functionNoExport(size_t i) {
  return getRuleContext<bapelParser::FunctionNoExportContext>(i);
}

bapelParser::InherentImplContext::InherentImplContext(ImplBlockContext *ctx) { copyFrom(ctx); }


antlrcpp::Any bapelParser::InherentImplContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitInherentImpl(this);
  else
    return visitor->visitChildren(this);
}
//----------------- TraitImplContext ------------------------------------------------------------------

tree::TerminalNode* bapelParser::TraitImplContext::IMPL() {
  return getToken(bapelParser::IMPL, 0);
}

std::vector<bapelParser::Type_Context *> bapelParser::TraitImplContext::type_() {
  return getRuleContexts<bapelParser::Type_Context>();
}

bapelParser::Type_Context* bapelParser::TraitImplContext::type_(size_t i) {
  return getRuleContext<bapelParser::Type_Context>(i);
}

tree::TerminalNode* bapelParser::TraitImplContext::FOR() {
  return getToken(bapelParser::FOR, 0);
}

tree::TerminalNode* bapelParser::TraitImplContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

tree::TerminalNode* bapelParser::TraitImplContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}

bapelParser::TypeAbstractionContext* bapelParser::TraitImplContext::typeAbstraction() {
  return getRuleContext<bapelParser::TypeAbstractionContext>(0);
}

std::vector<bapelParser::FunctionNoExportContext *> bapelParser::TraitImplContext::functionNoExport() {
  return getRuleContexts<bapelParser::FunctionNoExportContext>();
}

bapelParser::FunctionNoExportContext* bapelParser::TraitImplContext::functionNoExport(size_t i) {
  return getRuleContext<bapelParser::FunctionNoExportContext>(i);
}

bapelParser::TraitImplContext::TraitImplContext(ImplBlockContext *ctx) { copyFrom(ctx); }


antlrcpp::Any bapelParser::TraitImplContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTraitImpl(this);
  else
    return visitor->visitChildren(this);
}
bapelParser::ImplBlockContext* bapelParser::implBlock() {
  ImplBlockContext *_localctx = _tracker.createInstance<ImplBlockContext>(_ctx, getState());
  enterRule(_localctx, 30, bapelParser::RuleImplBlock);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(338);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 23, _ctx)) {
    case 1: {
      _localctx = dynamic_cast<ImplBlockContext *>(_tracker.createInstance<bapelParser::TraitImplContext>(_localctx));
      enterOuterAlt(_localctx, 1);
      setState(308);
      match(bapelParser::IMPL);
      setState(310);
      _errHandler->sync(this);

      switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 19, _ctx)) {
      case 1: {
        setState(309);
        typeAbstraction();
        break;
      }

      default:
        break;
      }
      setState(312);
      type_();
      setState(313);
      match(bapelParser::FOR);
      setState(314);
      type_();
      setState(315);
      match(bapelParser::LBRACE);
      setState(319);
      _errHandler->sync(this);
      _la = _input->LA(1);
      while (_la == bapelParser::FN) {
        setState(316);
        functionNoExport();
        setState(321);
        _errHandler->sync(this);
        _la = _input->LA(1);
      }
      setState(322);
      match(bapelParser::RBRACE);
      break;
    }

    case 2: {
      _localctx = dynamic_cast<ImplBlockContext *>(_tracker.createInstance<bapelParser::InherentImplContext>(_localctx));
      enterOuterAlt(_localctx, 2);
      setState(324);
      match(bapelParser::IMPL);
      setState(326);
      _errHandler->sync(this);

      switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 21, _ctx)) {
      case 1: {
        setState(325);
        typeAbstraction();
        break;
      }

      default:
        break;
      }
      setState(328);
      type_();
      setState(329);
      match(bapelParser::LBRACE);
      setState(333);
      _errHandler->sync(this);
      _la = _input->LA(1);
      while (_la == bapelParser::FN) {
        setState(330);
        functionNoExport();
        setState(335);
        _errHandler->sync(this);
        _la = _input->LA(1);
      }
      setState(336);
      match(bapelParser::RBRACE);
      break;
    }

    default:
      break;
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- DeclNoExportContext ------------------------------------------------------------------

bapelParser::DeclNoExportContext::DeclNoExportContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::DeclNoTermContext* bapelParser::DeclNoExportContext::declNoTerm() {
  return getRuleContext<bapelParser::DeclNoTermContext>(0);
}

tree::TerminalNode* bapelParser::DeclNoExportContext::DECL() {
  return getToken(bapelParser::DECL, 0);
}

bapelParser::TermDeclContext* bapelParser::DeclNoExportContext::termDecl() {
  return getRuleContext<bapelParser::TermDeclContext>(0);
}

tree::TerminalNode* bapelParser::DeclNoExportContext::PUB() {
  return getToken(bapelParser::PUB, 0);
}


size_t bapelParser::DeclNoExportContext::getRuleIndex() const {
  return bapelParser::RuleDeclNoExport;
}


antlrcpp::Any bapelParser::DeclNoExportContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitDeclNoExport(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::DeclNoExportContext* bapelParser::declNoExport() {
  DeclNoExportContext *_localctx = _tracker.createInstance<DeclNoExportContext>(_ctx, getState());
  enterRule(_localctx, 32, bapelParser::RuleDeclNoExport);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(345);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 24, _ctx)) {
    case 1: {
      enterOuterAlt(_localctx, 1);
      setState(340);
      declNoTerm();
      break;
    }

    case 2: {
      enterOuterAlt(_localctx, 2);
      setState(341);
      match(bapelParser::DECL);
      setState(342);
      termDecl();
      break;
    }

    case 3: {
      enterOuterAlt(_localctx, 3);
      setState(343);
      match(bapelParser::PUB);
      setState(344);
      termDecl();
      break;
    }

    default:
      break;
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- FunctionNoExportContext ------------------------------------------------------------------

bapelParser::FunctionNoExportContext::FunctionNoExportContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::FunctionNoExportContext::FN() {
  return getToken(bapelParser::FN, 0);
}

bapelParser::IdContext* bapelParser::FunctionNoExportContext::id() {
  return getRuleContext<bapelParser::IdContext>(0);
}

bapelParser::FunctionArgsContext* bapelParser::FunctionNoExportContext::functionArgs() {
  return getRuleContext<bapelParser::FunctionArgsContext>(0);
}

tree::TerminalNode* bapelParser::FunctionNoExportContext::ARROW() {
  return getToken(bapelParser::ARROW, 0);
}

bapelParser::Type_Context* bapelParser::FunctionNoExportContext::type_() {
  return getRuleContext<bapelParser::Type_Context>(0);
}

bapelParser::BlockExprContext* bapelParser::FunctionNoExportContext::blockExpr() {
  return getRuleContext<bapelParser::BlockExprContext>(0);
}

bapelParser::TypeAbstractionContext* bapelParser::FunctionNoExportContext::typeAbstraction() {
  return getRuleContext<bapelParser::TypeAbstractionContext>(0);
}


size_t bapelParser::FunctionNoExportContext::getRuleIndex() const {
  return bapelParser::RuleFunctionNoExport;
}


antlrcpp::Any bapelParser::FunctionNoExportContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitFunctionNoExport(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::FunctionNoExportContext* bapelParser::functionNoExport() {
  FunctionNoExportContext *_localctx = _tracker.createInstance<FunctionNoExportContext>(_ctx, getState());
  enterRule(_localctx, 34, bapelParser::RuleFunctionNoExport);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(347);
    match(bapelParser::FN);
    setState(348);
    id();
    setState(350);
    _errHandler->sync(this);

    _la = _input->LA(1);
    if (_la == bapelParser::LBRACKET) {
      setState(349);
      typeAbstraction();
    }
    setState(352);
    functionArgs();
    setState(353);
    match(bapelParser::ARROW);
    setState(354);
    type_();
    setState(355);
    blockExpr();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- FunctionArgsContext ------------------------------------------------------------------

bapelParser::FunctionArgsContext::FunctionArgsContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::FunctionArgsContext::LPAREN() {
  return getToken(bapelParser::LPAREN, 0);
}

tree::TerminalNode* bapelParser::FunctionArgsContext::RPAREN() {
  return getToken(bapelParser::RPAREN, 0);
}

std::vector<bapelParser::ArgContext *> bapelParser::FunctionArgsContext::arg() {
  return getRuleContexts<bapelParser::ArgContext>();
}

bapelParser::ArgContext* bapelParser::FunctionArgsContext::arg(size_t i) {
  return getRuleContext<bapelParser::ArgContext>(i);
}

std::vector<tree::TerminalNode *> bapelParser::FunctionArgsContext::COMMA() {
  return getTokens(bapelParser::COMMA);
}

tree::TerminalNode* bapelParser::FunctionArgsContext::COMMA(size_t i) {
  return getToken(bapelParser::COMMA, i);
}


size_t bapelParser::FunctionArgsContext::getRuleIndex() const {
  return bapelParser::RuleFunctionArgs;
}


antlrcpp::Any bapelParser::FunctionArgsContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitFunctionArgs(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::FunctionArgsContext* bapelParser::functionArgs() {
  FunctionArgsContext *_localctx = _tracker.createInstance<FunctionArgsContext>(_ctx, getState());
  enterRule(_localctx, 36, bapelParser::RuleFunctionArgs);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(357);
    match(bapelParser::LPAREN);
    setState(366);
    _errHandler->sync(this);

    _la = _input->LA(1);
    if (_la == bapelParser::IDENTIFIER) {
      setState(358);
      arg();
      setState(363);
      _errHandler->sync(this);
      _la = _input->LA(1);
      while (_la == bapelParser::COMMA) {
        setState(359);
        match(bapelParser::COMMA);
        setState(360);
        arg();
        setState(365);
        _errHandler->sync(this);
        _la = _input->LA(1);
      }
    }
    setState(368);
    match(bapelParser::RPAREN);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ArgContext ------------------------------------------------------------------

bapelParser::ArgContext::ArgContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::ArgContext::IDENTIFIER() {
  return getToken(bapelParser::IDENTIFIER, 0);
}

tree::TerminalNode* bapelParser::ArgContext::COLON() {
  return getToken(bapelParser::COLON, 0);
}

bapelParser::Type_Context* bapelParser::ArgContext::type_() {
  return getRuleContext<bapelParser::Type_Context>(0);
}


size_t bapelParser::ArgContext::getRuleIndex() const {
  return bapelParser::RuleArg;
}


antlrcpp::Any bapelParser::ArgContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitArg(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::ArgContext* bapelParser::arg() {
  ArgContext *_localctx = _tracker.createInstance<ArgContext>(_ctx, getState());
  enterRule(_localctx, 38, bapelParser::RuleArg);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(370);
    match(bapelParser::IDENTIFIER);
    setState(371);
    match(bapelParser::COLON);
    setState(372);
    type_();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- DeclContext ------------------------------------------------------------------

bapelParser::DeclContext::DeclContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::DeclContext::PUB() {
  return getToken(bapelParser::PUB, 0);
}

bapelParser::UnexportedDeclContext* bapelParser::DeclContext::unexportedDecl() {
  return getRuleContext<bapelParser::UnexportedDeclContext>(0);
}


size_t bapelParser::DeclContext::getRuleIndex() const {
  return bapelParser::RuleDecl;
}


antlrcpp::Any bapelParser::DeclContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitDecl(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::DeclContext* bapelParser::decl() {
  DeclContext *_localctx = _tracker.createInstance<DeclContext>(_ctx, getState());
  enterRule(_localctx, 40, bapelParser::RuleDecl);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(377);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::PUB: {
        enterOuterAlt(_localctx, 1);
        setState(374);
        match(bapelParser::PUB);
        setState(375);
        unexportedDecl();
        break;
      }

      case bapelParser::TYPE:
      case bapelParser::LPAREN:
      case bapelParser::IDENTIFIER: {
        enterOuterAlt(_localctx, 2);
        setState(376);
        unexportedDecl();
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- UnexportedDeclContext ------------------------------------------------------------------

bapelParser::UnexportedDeclContext::UnexportedDeclContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::TermDeclContext* bapelParser::UnexportedDeclContext::termDecl() {
  return getRuleContext<bapelParser::TermDeclContext>(0);
}

bapelParser::TypeDeclContext* bapelParser::UnexportedDeclContext::typeDecl() {
  return getRuleContext<bapelParser::TypeDeclContext>(0);
}


size_t bapelParser::UnexportedDeclContext::getRuleIndex() const {
  return bapelParser::RuleUnexportedDecl;
}


antlrcpp::Any bapelParser::UnexportedDeclContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitUnexportedDecl(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::UnexportedDeclContext* bapelParser::unexportedDecl() {
  UnexportedDeclContext *_localctx = _tracker.createInstance<UnexportedDeclContext>(_ctx, getState());
  enterRule(_localctx, 42, bapelParser::RuleUnexportedDecl);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(381);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::LPAREN:
      case bapelParser::IDENTIFIER: {
        enterOuterAlt(_localctx, 1);
        setState(379);
        termDecl();
        break;
      }

      case bapelParser::TYPE: {
        enterOuterAlt(_localctx, 2);
        setState(380);
        typeDecl();
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- DeclNoTermContext ------------------------------------------------------------------

bapelParser::DeclNoTermContext::DeclNoTermContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::DeclNoTermContext::PUB() {
  return getToken(bapelParser::PUB, 0);
}

bapelParser::TypeDeclContext* bapelParser::DeclNoTermContext::typeDecl() {
  return getRuleContext<bapelParser::TypeDeclContext>(0);
}


size_t bapelParser::DeclNoTermContext::getRuleIndex() const {
  return bapelParser::RuleDeclNoTerm;
}


antlrcpp::Any bapelParser::DeclNoTermContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitDeclNoTerm(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::DeclNoTermContext* bapelParser::declNoTerm() {
  DeclNoTermContext *_localctx = _tracker.createInstance<DeclNoTermContext>(_ctx, getState());
  enterRule(_localctx, 44, bapelParser::RuleDeclNoTerm);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(386);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::PUB: {
        enterOuterAlt(_localctx, 1);
        setState(383);
        match(bapelParser::PUB);
        setState(384);
        typeDecl();
        break;
      }

      case bapelParser::TYPE: {
        enterOuterAlt(_localctx, 2);
        setState(385);
        typeDecl();
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- TermDeclContext ------------------------------------------------------------------

bapelParser::TermDeclContext::TermDeclContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::IdContext* bapelParser::TermDeclContext::id() {
  return getRuleContext<bapelParser::IdContext>(0);
}

tree::TerminalNode* bapelParser::TermDeclContext::COLON() {
  return getToken(bapelParser::COLON, 0);
}

bapelParser::Type_Context* bapelParser::TermDeclContext::type_() {
  return getRuleContext<bapelParser::Type_Context>(0);
}


size_t bapelParser::TermDeclContext::getRuleIndex() const {
  return bapelParser::RuleTermDecl;
}


antlrcpp::Any bapelParser::TermDeclContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTermDecl(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TermDeclContext* bapelParser::termDecl() {
  TermDeclContext *_localctx = _tracker.createInstance<TermDeclContext>(_ctx, getState());
  enterRule(_localctx, 46, bapelParser::RuleTermDecl);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(388);
    id();
    setState(389);
    match(bapelParser::COLON);
    setState(390);
    type_();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- TypeDeclContext ------------------------------------------------------------------

bapelParser::TypeDeclContext::TypeDeclContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::TypeDeclContext::TYPE() {
  return getToken(bapelParser::TYPE, 0);
}

bapelParser::IdContext* bapelParser::TypeDeclContext::id() {
  return getRuleContext<bapelParser::IdContext>(0);
}

tree::TerminalNode* bapelParser::TypeDeclContext::ASSIGN() {
  return getToken(bapelParser::ASSIGN, 0);
}

bapelParser::Type_Context* bapelParser::TypeDeclContext::type_() {
  return getRuleContext<bapelParser::Type_Context>(0);
}

bapelParser::TypeAbstractionContext* bapelParser::TypeDeclContext::typeAbstraction() {
  return getRuleContext<bapelParser::TypeAbstractionContext>(0);
}


size_t bapelParser::TypeDeclContext::getRuleIndex() const {
  return bapelParser::RuleTypeDecl;
}


antlrcpp::Any bapelParser::TypeDeclContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTypeDecl(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TypeDeclContext* bapelParser::typeDecl() {
  TypeDeclContext *_localctx = _tracker.createInstance<TypeDeclContext>(_ctx, getState());
  enterRule(_localctx, 48, bapelParser::RuleTypeDecl);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(405);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 33, _ctx)) {
    case 1: {
      enterOuterAlt(_localctx, 1);
      setState(392);
      match(bapelParser::TYPE);
      setState(393);
      id();
      setState(395);
      _errHandler->sync(this);

      _la = _input->LA(1);
      if (_la == bapelParser::LBRACKET) {
        setState(394);
        typeAbstraction();
      }
      setState(397);
      match(bapelParser::ASSIGN);
      setState(398);
      type_();
      break;
    }

    case 2: {
      enterOuterAlt(_localctx, 2);
      setState(400);
      match(bapelParser::TYPE);
      setState(401);
      id();
      setState(403);
      _errHandler->sync(this);

      _la = _input->LA(1);
      if (_la == bapelParser::LBRACKET) {
        setState(402);
        typeAbstraction();
      }
      break;
    }

    default:
      break;
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- TypeAbstractionContext ------------------------------------------------------------------

bapelParser::TypeAbstractionContext::TypeAbstractionContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::TypeAbstractionContext::LBRACKET() {
  return getToken(bapelParser::LBRACKET, 0);
}

std::vector<bapelParser::BoundedTvarContext *> bapelParser::TypeAbstractionContext::boundedTvar() {
  return getRuleContexts<bapelParser::BoundedTvarContext>();
}

bapelParser::BoundedTvarContext* bapelParser::TypeAbstractionContext::boundedTvar(size_t i) {
  return getRuleContext<bapelParser::BoundedTvarContext>(i);
}

tree::TerminalNode* bapelParser::TypeAbstractionContext::RBRACKET() {
  return getToken(bapelParser::RBRACKET, 0);
}

std::vector<tree::TerminalNode *> bapelParser::TypeAbstractionContext::COMMA() {
  return getTokens(bapelParser::COMMA);
}

tree::TerminalNode* bapelParser::TypeAbstractionContext::COMMA(size_t i) {
  return getToken(bapelParser::COMMA, i);
}


size_t bapelParser::TypeAbstractionContext::getRuleIndex() const {
  return bapelParser::RuleTypeAbstraction;
}


antlrcpp::Any bapelParser::TypeAbstractionContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTypeAbstraction(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TypeAbstractionContext* bapelParser::typeAbstraction() {
  TypeAbstractionContext *_localctx = _tracker.createInstance<TypeAbstractionContext>(_ctx, getState());
  enterRule(_localctx, 50, bapelParser::RuleTypeAbstraction);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(407);
    match(bapelParser::LBRACKET);
    setState(408);
    boundedTvar();
    setState(413);
    _errHandler->sync(this);
    _la = _input->LA(1);
    while (_la == bapelParser::COMMA) {
      setState(409);
      match(bapelParser::COMMA);
      setState(410);
      boundedTvar();
      setState(415);
      _errHandler->sync(this);
      _la = _input->LA(1);
    }
    setState(416);
    match(bapelParser::RBRACKET);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- BoundedTvarContext ------------------------------------------------------------------

bapelParser::BoundedTvarContext::BoundedTvarContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::TvarContext* bapelParser::BoundedTvarContext::tvar() {
  return getRuleContext<bapelParser::TvarContext>(0);
}

tree::TerminalNode* bapelParser::BoundedTvarContext::COLON() {
  return getToken(bapelParser::COLON, 0);
}

bapelParser::TraitBoundContext* bapelParser::BoundedTvarContext::traitBound() {
  return getRuleContext<bapelParser::TraitBoundContext>(0);
}


size_t bapelParser::BoundedTvarContext::getRuleIndex() const {
  return bapelParser::RuleBoundedTvar;
}


antlrcpp::Any bapelParser::BoundedTvarContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitBoundedTvar(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::BoundedTvarContext* bapelParser::boundedTvar() {
  BoundedTvarContext *_localctx = _tracker.createInstance<BoundedTvarContext>(_ctx, getState());
  enterRule(_localctx, 52, bapelParser::RuleBoundedTvar);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(418);
    tvar();
    setState(421);
    _errHandler->sync(this);

    _la = _input->LA(1);
    if (_la == bapelParser::COLON) {
      setState(419);
      match(bapelParser::COLON);
      setState(420);
      traitBound();
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- TvarContext ------------------------------------------------------------------

bapelParser::TvarContext::TvarContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::TvarContext::SINGLE_QUOTE() {
  return getToken(bapelParser::SINGLE_QUOTE, 0);
}

tree::TerminalNode* bapelParser::TvarContext::IDENTIFIER() {
  return getToken(bapelParser::IDENTIFIER, 0);
}


size_t bapelParser::TvarContext::getRuleIndex() const {
  return bapelParser::RuleTvar;
}


antlrcpp::Any bapelParser::TvarContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTvar(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TvarContext* bapelParser::tvar() {
  TvarContext *_localctx = _tracker.createInstance<TvarContext>(_ctx, getState());
  enterRule(_localctx, 54, bapelParser::RuleTvar);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(423);
    match(bapelParser::SINGLE_QUOTE);
    setState(424);
    match(bapelParser::IDENTIFIER);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- TraitBoundContext ------------------------------------------------------------------

bapelParser::TraitBoundContext::TraitBoundContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

std::vector<bapelParser::Type_Context *> bapelParser::TraitBoundContext::type_() {
  return getRuleContexts<bapelParser::Type_Context>();
}

bapelParser::Type_Context* bapelParser::TraitBoundContext::type_(size_t i) {
  return getRuleContext<bapelParser::Type_Context>(i);
}

std::vector<tree::TerminalNode *> bapelParser::TraitBoundContext::PLUS() {
  return getTokens(bapelParser::PLUS);
}

tree::TerminalNode* bapelParser::TraitBoundContext::PLUS(size_t i) {
  return getToken(bapelParser::PLUS, i);
}


size_t bapelParser::TraitBoundContext::getRuleIndex() const {
  return bapelParser::RuleTraitBound;
}


antlrcpp::Any bapelParser::TraitBoundContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTraitBound(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TraitBoundContext* bapelParser::traitBound() {
  TraitBoundContext *_localctx = _tracker.createInstance<TraitBoundContext>(_ctx, getState());
  enterRule(_localctx, 56, bapelParser::RuleTraitBound);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(426);
    type_();
    setState(431);
    _errHandler->sync(this);
    _la = _input->LA(1);
    while (_la == bapelParser::PLUS) {
      setState(427);
      match(bapelParser::PLUS);
      setState(428);
      type_();
      setState(433);
      _errHandler->sync(this);
      _la = _input->LA(1);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- Type_Context ------------------------------------------------------------------

bapelParser::Type_Context::Type_Context(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::ForallTypeContext* bapelParser::Type_Context::forallType() {
  return getRuleContext<bapelParser::ForallTypeContext>(0);
}


size_t bapelParser::Type_Context::getRuleIndex() const {
  return bapelParser::RuleType_;
}


antlrcpp::Any bapelParser::Type_Context::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitType_(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::Type_Context* bapelParser::type_() {
  Type_Context *_localctx = _tracker.createInstance<Type_Context>(_ctx, getState());
  enterRule(_localctx, 58, bapelParser::RuleType_);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(434);
    forallType();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ForallTypeContext ------------------------------------------------------------------

bapelParser::ForallTypeContext::ForallTypeContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::ForallTypeContext::FORALL() {
  return getToken(bapelParser::FORALL, 0);
}

bapelParser::TypeAbstractionContext* bapelParser::ForallTypeContext::typeAbstraction() {
  return getRuleContext<bapelParser::TypeAbstractionContext>(0);
}

bapelParser::FunctionTypeContext* bapelParser::ForallTypeContext::functionType() {
  return getRuleContext<bapelParser::FunctionTypeContext>(0);
}


size_t bapelParser::ForallTypeContext::getRuleIndex() const {
  return bapelParser::RuleForallType;
}


antlrcpp::Any bapelParser::ForallTypeContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitForallType(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::ForallTypeContext* bapelParser::forallType() {
  ForallTypeContext *_localctx = _tracker.createInstance<ForallTypeContext>(_ctx, getState());
  enterRule(_localctx, 60, bapelParser::RuleForallType);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(441);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::FORALL: {
        enterOuterAlt(_localctx, 1);
        setState(436);
        match(bapelParser::FORALL);
        setState(437);
        typeAbstraction();
        setState(438);
        functionType();
        break;
      }

      case bapelParser::STRUCT:
      case bapelParser::VARIANT:
      case bapelParser::AMP:
      case bapelParser::LPAREN:
      case bapelParser::LBRACKET:
      case bapelParser::SINGLE_QUOTE:
      case bapelParser::IDENTIFIER: {
        enterOuterAlt(_localctx, 2);
        setState(440);
        functionType();
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- FunctionTypeContext ------------------------------------------------------------------

bapelParser::FunctionTypeContext::FunctionTypeContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::PtrTypeContext* bapelParser::FunctionTypeContext::ptrType() {
  return getRuleContext<bapelParser::PtrTypeContext>(0);
}

tree::TerminalNode* bapelParser::FunctionTypeContext::ARROW() {
  return getToken(bapelParser::ARROW, 0);
}

bapelParser::FunctionTypeContext* bapelParser::FunctionTypeContext::functionType() {
  return getRuleContext<bapelParser::FunctionTypeContext>(0);
}


size_t bapelParser::FunctionTypeContext::getRuleIndex() const {
  return bapelParser::RuleFunctionType;
}


antlrcpp::Any bapelParser::FunctionTypeContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitFunctionType(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::FunctionTypeContext* bapelParser::functionType() {
  FunctionTypeContext *_localctx = _tracker.createInstance<FunctionTypeContext>(_ctx, getState());
  enterRule(_localctx, 62, bapelParser::RuleFunctionType);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(443);
    ptrType();
    setState(446);
    _errHandler->sync(this);

    _la = _input->LA(1);
    if (_la == bapelParser::ARROW) {
      setState(444);
      match(bapelParser::ARROW);
      setState(445);
      functionType();
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- PtrTypeContext ------------------------------------------------------------------

bapelParser::PtrTypeContext::PtrTypeContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::PtrTypeContext::AMP() {
  return getToken(bapelParser::AMP, 0);
}

bapelParser::PtrTypeContext* bapelParser::PtrTypeContext::ptrType() {
  return getRuleContext<bapelParser::PtrTypeContext>(0);
}

bapelParser::AppTypeContext* bapelParser::PtrTypeContext::appType() {
  return getRuleContext<bapelParser::AppTypeContext>(0);
}


size_t bapelParser::PtrTypeContext::getRuleIndex() const {
  return bapelParser::RulePtrType;
}


antlrcpp::Any bapelParser::PtrTypeContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitPtrType(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::PtrTypeContext* bapelParser::ptrType() {
  PtrTypeContext *_localctx = _tracker.createInstance<PtrTypeContext>(_ctx, getState());
  enterRule(_localctx, 64, bapelParser::RulePtrType);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(451);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::AMP: {
        enterOuterAlt(_localctx, 1);
        setState(448);
        match(bapelParser::AMP);
        setState(449);
        ptrType();
        break;
      }

      case bapelParser::STRUCT:
      case bapelParser::VARIANT:
      case bapelParser::LPAREN:
      case bapelParser::LBRACKET:
      case bapelParser::SINGLE_QUOTE:
      case bapelParser::IDENTIFIER: {
        enterOuterAlt(_localctx, 2);
        setState(450);
        appType(0);
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- AppTypeContext ------------------------------------------------------------------

bapelParser::AppTypeContext::AppTypeContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::PrimaryTypeContext* bapelParser::AppTypeContext::primaryType() {
  return getRuleContext<bapelParser::PrimaryTypeContext>(0);
}

bapelParser::AppTypeContext* bapelParser::AppTypeContext::appType() {
  return getRuleContext<bapelParser::AppTypeContext>(0);
}


size_t bapelParser::AppTypeContext::getRuleIndex() const {
  return bapelParser::RuleAppType;
}


antlrcpp::Any bapelParser::AppTypeContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitAppType(this);
  else
    return visitor->visitChildren(this);
}


bapelParser::AppTypeContext* bapelParser::appType() {
   return appType(0);
}

bapelParser::AppTypeContext* bapelParser::appType(int precedence) {
  ParserRuleContext *parentContext = _ctx;
  size_t parentState = getState();
  bapelParser::AppTypeContext *_localctx = _tracker.createInstance<AppTypeContext>(_ctx, parentState);
  bapelParser::AppTypeContext *previousContext = _localctx;
  (void)previousContext; // Silence compiler, in case the context is not used by generated code.
  size_t startState = 66;
  enterRecursionRule(_localctx, 66, bapelParser::RuleAppType, precedence);

    

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    unrollRecursionContexts(parentContext);
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(454);
    primaryType();
    _ctx->stop = _input->LT(-1);
    setState(460);
    _errHandler->sync(this);
    alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 40, _ctx);
    while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER) {
      if (alt == 1) {
        if (!_parseListeners.empty())
          triggerExitRuleEvent();
        previousContext = _localctx;
        _localctx = _tracker.createInstance<AppTypeContext>(parentContext, parentState);
        pushNewRecursionContext(_localctx, startState, RuleAppType);
        setState(456);

        if (!(precpred(_ctx, 2))) throw FailedPredicateException(this, "precpred(_ctx, 2)");
        setState(457);
        primaryType(); 
      }
      setState(462);
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 40, _ctx);
    }
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }
  return _localctx;
}

//----------------- PrimaryTypeContext ------------------------------------------------------------------

bapelParser::PrimaryTypeContext::PrimaryTypeContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::ArrayTypeContext* bapelParser::PrimaryTypeContext::arrayType() {
  return getRuleContext<bapelParser::ArrayTypeContext>(0);
}

bapelParser::StructTypeContext* bapelParser::PrimaryTypeContext::structType() {
  return getRuleContext<bapelParser::StructTypeContext>(0);
}

bapelParser::TupleTypeContext* bapelParser::PrimaryTypeContext::tupleType() {
  return getRuleContext<bapelParser::TupleTypeContext>(0);
}

bapelParser::VariantTypeContext* bapelParser::PrimaryTypeContext::variantType() {
  return getRuleContext<bapelParser::VariantTypeContext>(0);
}

tree::TerminalNode* bapelParser::PrimaryTypeContext::SINGLE_QUOTE() {
  return getToken(bapelParser::SINGLE_QUOTE, 0);
}

tree::TerminalNode* bapelParser::PrimaryTypeContext::IDENTIFIER() {
  return getToken(bapelParser::IDENTIFIER, 0);
}

bapelParser::IdContext* bapelParser::PrimaryTypeContext::id() {
  return getRuleContext<bapelParser::IdContext>(0);
}

tree::TerminalNode* bapelParser::PrimaryTypeContext::LPAREN() {
  return getToken(bapelParser::LPAREN, 0);
}

bapelParser::Type_Context* bapelParser::PrimaryTypeContext::type_() {
  return getRuleContext<bapelParser::Type_Context>(0);
}

tree::TerminalNode* bapelParser::PrimaryTypeContext::RPAREN() {
  return getToken(bapelParser::RPAREN, 0);
}


size_t bapelParser::PrimaryTypeContext::getRuleIndex() const {
  return bapelParser::RulePrimaryType;
}


antlrcpp::Any bapelParser::PrimaryTypeContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitPrimaryType(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::PrimaryTypeContext* bapelParser::primaryType() {
  PrimaryTypeContext *_localctx = _tracker.createInstance<PrimaryTypeContext>(_ctx, getState());
  enterRule(_localctx, 68, bapelParser::RulePrimaryType);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(474);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 41, _ctx)) {
    case 1: {
      enterOuterAlt(_localctx, 1);
      setState(463);
      arrayType();
      break;
    }

    case 2: {
      enterOuterAlt(_localctx, 2);
      setState(464);
      structType();
      break;
    }

    case 3: {
      enterOuterAlt(_localctx, 3);
      setState(465);
      tupleType();
      break;
    }

    case 4: {
      enterOuterAlt(_localctx, 4);
      setState(466);
      variantType();
      break;
    }

    case 5: {
      enterOuterAlt(_localctx, 5);
      setState(467);
      match(bapelParser::SINGLE_QUOTE);
      setState(468);
      match(bapelParser::IDENTIFIER);
      break;
    }

    case 6: {
      enterOuterAlt(_localctx, 6);
      setState(469);
      id();
      break;
    }

    case 7: {
      enterOuterAlt(_localctx, 7);
      setState(470);
      match(bapelParser::LPAREN);
      setState(471);
      type_();
      setState(472);
      match(bapelParser::RPAREN);
      break;
    }

    default:
      break;
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ArrayTypeContext ------------------------------------------------------------------

bapelParser::ArrayTypeContext::ArrayTypeContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::ArrayTypeContext::LBRACKET() {
  return getToken(bapelParser::LBRACKET, 0);
}

bapelParser::Type_Context* bapelParser::ArrayTypeContext::type_() {
  return getRuleContext<bapelParser::Type_Context>(0);
}

tree::TerminalNode* bapelParser::ArrayTypeContext::COMMA() {
  return getToken(bapelParser::COMMA, 0);
}

tree::TerminalNode* bapelParser::ArrayTypeContext::INT_LITERAL() {
  return getToken(bapelParser::INT_LITERAL, 0);
}

tree::TerminalNode* bapelParser::ArrayTypeContext::RBRACKET() {
  return getToken(bapelParser::RBRACKET, 0);
}

tree::TerminalNode* bapelParser::ArrayTypeContext::MINUS() {
  return getToken(bapelParser::MINUS, 0);
}


size_t bapelParser::ArrayTypeContext::getRuleIndex() const {
  return bapelParser::RuleArrayType;
}


antlrcpp::Any bapelParser::ArrayTypeContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitArrayType(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::ArrayTypeContext* bapelParser::arrayType() {
  ArrayTypeContext *_localctx = _tracker.createInstance<ArrayTypeContext>(_ctx, getState());
  enterRule(_localctx, 70, bapelParser::RuleArrayType);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(493);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 42, _ctx)) {
    case 1: {
      enterOuterAlt(_localctx, 1);
      setState(476);
      match(bapelParser::LBRACKET);
      setState(477);
      type_();
      setState(478);
      match(bapelParser::COMMA);
      setState(479);
      match(bapelParser::INT_LITERAL);
      setState(480);
      match(bapelParser::RBRACKET);
      break;
    }

    case 2: {
      enterOuterAlt(_localctx, 2);
      setState(482);
      match(bapelParser::LBRACKET);
      setState(483);
      type_();
      setState(484);
      match(bapelParser::RBRACKET);
      break;
    }

    case 3: {
      enterOuterAlt(_localctx, 3);
      setState(486);
      match(bapelParser::LBRACKET);
      setState(487);
      type_();
      setState(488);
      match(bapelParser::COMMA);
      setState(489);
      match(bapelParser::MINUS);
      setState(490);
      match(bapelParser::INT_LITERAL);
      setState(491);
      match(bapelParser::RBRACKET);
      break;
    }

    default:
      break;
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- StructTypeContext ------------------------------------------------------------------

bapelParser::StructTypeContext::StructTypeContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::StructTypeContext::STRUCT() {
  return getToken(bapelParser::STRUCT, 0);
}

tree::TerminalNode* bapelParser::StructTypeContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

tree::TerminalNode* bapelParser::StructTypeContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}

bapelParser::FieldsContext* bapelParser::StructTypeContext::fields() {
  return getRuleContext<bapelParser::FieldsContext>(0);
}

tree::TerminalNode* bapelParser::StructTypeContext::COMMA() {
  return getToken(bapelParser::COMMA, 0);
}


size_t bapelParser::StructTypeContext::getRuleIndex() const {
  return bapelParser::RuleStructType;
}


antlrcpp::Any bapelParser::StructTypeContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitStructType(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::StructTypeContext* bapelParser::structType() {
  StructTypeContext *_localctx = _tracker.createInstance<StructTypeContext>(_ctx, getState());
  enterRule(_localctx, 72, bapelParser::RuleStructType);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(495);
    match(bapelParser::STRUCT);
    setState(496);
    match(bapelParser::LBRACE);
    setState(501);
    _errHandler->sync(this);

    _la = _input->LA(1);
    if (_la == bapelParser::LPAREN

    || _la == bapelParser::IDENTIFIER) {
      setState(497);
      fields();

      setState(499);
      _errHandler->sync(this);

      _la = _input->LA(1);
      if (_la == bapelParser::COMMA) {
        setState(498);
        match(bapelParser::COMMA);
      }
    }
    setState(503);
    match(bapelParser::RBRACE);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- FieldsContext ------------------------------------------------------------------

bapelParser::FieldsContext::FieldsContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

std::vector<bapelParser::FieldContext *> bapelParser::FieldsContext::field() {
  return getRuleContexts<bapelParser::FieldContext>();
}

bapelParser::FieldContext* bapelParser::FieldsContext::field(size_t i) {
  return getRuleContext<bapelParser::FieldContext>(i);
}

std::vector<tree::TerminalNode *> bapelParser::FieldsContext::COMMA() {
  return getTokens(bapelParser::COMMA);
}

tree::TerminalNode* bapelParser::FieldsContext::COMMA(size_t i) {
  return getToken(bapelParser::COMMA, i);
}


size_t bapelParser::FieldsContext::getRuleIndex() const {
  return bapelParser::RuleFields;
}


antlrcpp::Any bapelParser::FieldsContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitFields(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::FieldsContext* bapelParser::fields() {
  FieldsContext *_localctx = _tracker.createInstance<FieldsContext>(_ctx, getState());
  enterRule(_localctx, 74, bapelParser::RuleFields);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(505);
    field();
    setState(510);
    _errHandler->sync(this);
    alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 45, _ctx);
    while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER) {
      if (alt == 1) {
        setState(506);
        match(bapelParser::COMMA);
        setState(507);
        field(); 
      }
      setState(512);
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 45, _ctx);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- FieldContext ------------------------------------------------------------------

bapelParser::FieldContext::FieldContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::IdContext* bapelParser::FieldContext::id() {
  return getRuleContext<bapelParser::IdContext>(0);
}

tree::TerminalNode* bapelParser::FieldContext::COLON() {
  return getToken(bapelParser::COLON, 0);
}

bapelParser::Type_Context* bapelParser::FieldContext::type_() {
  return getRuleContext<bapelParser::Type_Context>(0);
}


size_t bapelParser::FieldContext::getRuleIndex() const {
  return bapelParser::RuleField;
}


antlrcpp::Any bapelParser::FieldContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitField(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::FieldContext* bapelParser::field() {
  FieldContext *_localctx = _tracker.createInstance<FieldContext>(_ctx, getState());
  enterRule(_localctx, 76, bapelParser::RuleField);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(513);
    id();
    setState(514);
    match(bapelParser::COLON);
    setState(515);
    type_();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- TupleTypeContext ------------------------------------------------------------------

bapelParser::TupleTypeContext::TupleTypeContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::TupleTypeContext::LPAREN() {
  return getToken(bapelParser::LPAREN, 0);
}

tree::TerminalNode* bapelParser::TupleTypeContext::RPAREN() {
  return getToken(bapelParser::RPAREN, 0);
}

bapelParser::TupleTypeArgsContext* bapelParser::TupleTypeContext::tupleTypeArgs() {
  return getRuleContext<bapelParser::TupleTypeArgsContext>(0);
}


size_t bapelParser::TupleTypeContext::getRuleIndex() const {
  return bapelParser::RuleTupleType;
}


antlrcpp::Any bapelParser::TupleTypeContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTupleType(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TupleTypeContext* bapelParser::tupleType() {
  TupleTypeContext *_localctx = _tracker.createInstance<TupleTypeContext>(_ctx, getState());
  enterRule(_localctx, 78, bapelParser::RuleTupleType);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(523);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 46, _ctx)) {
    case 1: {
      enterOuterAlt(_localctx, 1);
      setState(517);
      match(bapelParser::LPAREN);
      setState(518);
      match(bapelParser::RPAREN);
      break;
    }

    case 2: {
      enterOuterAlt(_localctx, 2);
      setState(519);
      match(bapelParser::LPAREN);
      setState(520);
      tupleTypeArgs();
      setState(521);
      match(bapelParser::RPAREN);
      break;
    }

    default:
      break;
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- TupleTypeArgsContext ------------------------------------------------------------------

bapelParser::TupleTypeArgsContext::TupleTypeArgsContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

std::vector<bapelParser::Type_Context *> bapelParser::TupleTypeArgsContext::type_() {
  return getRuleContexts<bapelParser::Type_Context>();
}

bapelParser::Type_Context* bapelParser::TupleTypeArgsContext::type_(size_t i) {
  return getRuleContext<bapelParser::Type_Context>(i);
}

std::vector<tree::TerminalNode *> bapelParser::TupleTypeArgsContext::COMMA() {
  return getTokens(bapelParser::COMMA);
}

tree::TerminalNode* bapelParser::TupleTypeArgsContext::COMMA(size_t i) {
  return getToken(bapelParser::COMMA, i);
}


size_t bapelParser::TupleTypeArgsContext::getRuleIndex() const {
  return bapelParser::RuleTupleTypeArgs;
}


antlrcpp::Any bapelParser::TupleTypeArgsContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTupleTypeArgs(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TupleTypeArgsContext* bapelParser::tupleTypeArgs() {
  TupleTypeArgsContext *_localctx = _tracker.createInstance<TupleTypeArgsContext>(_ctx, getState());
  enterRule(_localctx, 80, bapelParser::RuleTupleTypeArgs);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(525);
    type_();
    setState(528); 
    _errHandler->sync(this);
    _la = _input->LA(1);
    do {
      setState(526);
      match(bapelParser::COMMA);
      setState(527);
      type_();
      setState(530); 
      _errHandler->sync(this);
      _la = _input->LA(1);
    } while (_la == bapelParser::COMMA);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- VariantTypeContext ------------------------------------------------------------------

bapelParser::VariantTypeContext::VariantTypeContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::VariantTypeContext::VARIANT() {
  return getToken(bapelParser::VARIANT, 0);
}

tree::TerminalNode* bapelParser::VariantTypeContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

tree::TerminalNode* bapelParser::VariantTypeContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}

bapelParser::TagsContext* bapelParser::VariantTypeContext::tags() {
  return getRuleContext<bapelParser::TagsContext>(0);
}

tree::TerminalNode* bapelParser::VariantTypeContext::COMMA() {
  return getToken(bapelParser::COMMA, 0);
}


size_t bapelParser::VariantTypeContext::getRuleIndex() const {
  return bapelParser::RuleVariantType;
}


antlrcpp::Any bapelParser::VariantTypeContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitVariantType(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::VariantTypeContext* bapelParser::variantType() {
  VariantTypeContext *_localctx = _tracker.createInstance<VariantTypeContext>(_ctx, getState());
  enterRule(_localctx, 82, bapelParser::RuleVariantType);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(532);
    match(bapelParser::VARIANT);
    setState(533);
    match(bapelParser::LBRACE);
    setState(538);
    _errHandler->sync(this);

    _la = _input->LA(1);
    if (_la == bapelParser::LPAREN

    || _la == bapelParser::IDENTIFIER) {
      setState(534);
      tags();

      setState(536);
      _errHandler->sync(this);

      _la = _input->LA(1);
      if (_la == bapelParser::COMMA) {
        setState(535);
        match(bapelParser::COMMA);
      }
    }
    setState(540);
    match(bapelParser::RBRACE);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- TagsContext ------------------------------------------------------------------

bapelParser::TagsContext::TagsContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

std::vector<bapelParser::TagContext *> bapelParser::TagsContext::tag() {
  return getRuleContexts<bapelParser::TagContext>();
}

bapelParser::TagContext* bapelParser::TagsContext::tag(size_t i) {
  return getRuleContext<bapelParser::TagContext>(i);
}

std::vector<tree::TerminalNode *> bapelParser::TagsContext::COMMA() {
  return getTokens(bapelParser::COMMA);
}

tree::TerminalNode* bapelParser::TagsContext::COMMA(size_t i) {
  return getToken(bapelParser::COMMA, i);
}


size_t bapelParser::TagsContext::getRuleIndex() const {
  return bapelParser::RuleTags;
}


antlrcpp::Any bapelParser::TagsContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTags(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TagsContext* bapelParser::tags() {
  TagsContext *_localctx = _tracker.createInstance<TagsContext>(_ctx, getState());
  enterRule(_localctx, 84, bapelParser::RuleTags);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(542);
    tag();
    setState(547);
    _errHandler->sync(this);
    alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 50, _ctx);
    while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER) {
      if (alt == 1) {
        setState(543);
        match(bapelParser::COMMA);
        setState(544);
        tag(); 
      }
      setState(549);
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 50, _ctx);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- TagContext ------------------------------------------------------------------

bapelParser::TagContext::TagContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::IdContext* bapelParser::TagContext::id() {
  return getRuleContext<bapelParser::IdContext>(0);
}

bapelParser::Type_Context* bapelParser::TagContext::type_() {
  return getRuleContext<bapelParser::Type_Context>(0);
}


size_t bapelParser::TagContext::getRuleIndex() const {
  return bapelParser::RuleTag;
}


antlrcpp::Any bapelParser::TagContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTag(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TagContext* bapelParser::tag() {
  TagContext *_localctx = _tracker.createInstance<TagContext>(_ctx, getState());
  enterRule(_localctx, 86, bapelParser::RuleTag);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(550);
    id();
    setState(551);
    type_();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ExpressionContext ------------------------------------------------------------------

bapelParser::ExpressionContext::ExpressionContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::ExpressionWithoutBlockContext* bapelParser::ExpressionContext::expressionWithoutBlock() {
  return getRuleContext<bapelParser::ExpressionWithoutBlockContext>(0);
}

bapelParser::ExpressionWithBlockContext* bapelParser::ExpressionContext::expressionWithBlock() {
  return getRuleContext<bapelParser::ExpressionWithBlockContext>(0);
}


size_t bapelParser::ExpressionContext::getRuleIndex() const {
  return bapelParser::RuleExpression;
}


antlrcpp::Any bapelParser::ExpressionContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitExpression(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::ExpressionContext* bapelParser::expression() {
  ExpressionContext *_localctx = _tracker.createInstance<ExpressionContext>(_ctx, getState());
  enterRule(_localctx, 88, bapelParser::RuleExpression);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(555);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::STRUCT:
      case bapelParser::VARIANT:
      case bapelParser::RETURN:
      case bapelParser::MINUS:
      case bapelParser::MUL:
      case bapelParser::NOT:
      case bapelParser::AMP:
      case bapelParser::LPAREN:
      case bapelParser::IDENTIFIER:
      case bapelParser::INT_LITERAL:
      case bapelParser::FLOAT_LITERAL:
      case bapelParser::RUNE_LITERAL:
      case bapelParser::STRING_LITERAL: {
        enterOuterAlt(_localctx, 1);
        setState(553);
        expressionWithoutBlock();
        break;
      }

      case bapelParser::FN:
      case bapelParser::MATCH:
      case bapelParser::SET:
      case bapelParser::IF:
      case bapelParser::FOR:
      case bapelParser::LBRACE: {
        enterOuterAlt(_localctx, 2);
        setState(554);
        expressionWithBlock();
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ExpressionWithoutBlockContext ------------------------------------------------------------------

bapelParser::ExpressionWithoutBlockContext::ExpressionWithoutBlockContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::AssignTermContext* bapelParser::ExpressionWithoutBlockContext::assignTerm() {
  return getRuleContext<bapelParser::AssignTermContext>(0);
}

bapelParser::OperatorExprContext* bapelParser::ExpressionWithoutBlockContext::operatorExpr() {
  return getRuleContext<bapelParser::OperatorExprContext>(0);
}

bapelParser::ReturnTermContext* bapelParser::ExpressionWithoutBlockContext::returnTerm() {
  return getRuleContext<bapelParser::ReturnTermContext>(0);
}


size_t bapelParser::ExpressionWithoutBlockContext::getRuleIndex() const {
  return bapelParser::RuleExpressionWithoutBlock;
}


antlrcpp::Any bapelParser::ExpressionWithoutBlockContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitExpressionWithoutBlock(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::ExpressionWithoutBlockContext* bapelParser::expressionWithoutBlock() {
  ExpressionWithoutBlockContext *_localctx = _tracker.createInstance<ExpressionWithoutBlockContext>(_ctx, getState());
  enterRule(_localctx, 90, bapelParser::RuleExpressionWithoutBlock);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(560);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 52, _ctx)) {
    case 1: {
      enterOuterAlt(_localctx, 1);
      setState(557);
      assignTerm();
      break;
    }

    case 2: {
      enterOuterAlt(_localctx, 2);
      setState(558);
      operatorExpr();
      break;
    }

    case 3: {
      enterOuterAlt(_localctx, 3);
      setState(559);
      returnTerm();
      break;
    }

    default:
      break;
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ExpressionWithBlockContext ------------------------------------------------------------------

bapelParser::ExpressionWithBlockContext::ExpressionWithBlockContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::BlockExprContext* bapelParser::ExpressionWithBlockContext::blockExpr() {
  return getRuleContext<bapelParser::BlockExprContext>(0);
}

bapelParser::IfTermContext* bapelParser::ExpressionWithBlockContext::ifTerm() {
  return getRuleContext<bapelParser::IfTermContext>(0);
}

bapelParser::ForTermContext* bapelParser::ExpressionWithBlockContext::forTerm() {
  return getRuleContext<bapelParser::ForTermContext>(0);
}

bapelParser::LambdaTermContext* bapelParser::ExpressionWithBlockContext::lambdaTerm() {
  return getRuleContext<bapelParser::LambdaTermContext>(0);
}

bapelParser::MatchTermContext* bapelParser::ExpressionWithBlockContext::matchTerm() {
  return getRuleContext<bapelParser::MatchTermContext>(0);
}

bapelParser::SetTermContext* bapelParser::ExpressionWithBlockContext::setTerm() {
  return getRuleContext<bapelParser::SetTermContext>(0);
}


size_t bapelParser::ExpressionWithBlockContext::getRuleIndex() const {
  return bapelParser::RuleExpressionWithBlock;
}


antlrcpp::Any bapelParser::ExpressionWithBlockContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitExpressionWithBlock(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::ExpressionWithBlockContext* bapelParser::expressionWithBlock() {
  ExpressionWithBlockContext *_localctx = _tracker.createInstance<ExpressionWithBlockContext>(_ctx, getState());
  enterRule(_localctx, 92, bapelParser::RuleExpressionWithBlock);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(568);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::LBRACE: {
        enterOuterAlt(_localctx, 1);
        setState(562);
        blockExpr();
        break;
      }

      case bapelParser::IF: {
        enterOuterAlt(_localctx, 2);
        setState(563);
        ifTerm();
        break;
      }

      case bapelParser::FOR: {
        enterOuterAlt(_localctx, 3);
        setState(564);
        forTerm();
        break;
      }

      case bapelParser::FN: {
        enterOuterAlt(_localctx, 4);
        setState(565);
        lambdaTerm();
        break;
      }

      case bapelParser::MATCH: {
        enterOuterAlt(_localctx, 5);
        setState(566);
        matchTerm();
        break;
      }

      case bapelParser::SET: {
        enterOuterAlt(_localctx, 6);
        setState(567);
        setTerm();
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- AssignTermContext ------------------------------------------------------------------

bapelParser::AssignTermContext::AssignTermContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::AssignTermContext::LARROW() {
  return getToken(bapelParser::LARROW, 0);
}

bapelParser::ExpressionContext* bapelParser::AssignTermContext::expression() {
  return getRuleContext<bapelParser::ExpressionContext>(0);
}

bapelParser::IdContext* bapelParser::AssignTermContext::id() {
  return getRuleContext<bapelParser::IdContext>(0);
}

bapelParser::TupleExprContext* bapelParser::AssignTermContext::tupleExpr() {
  return getRuleContext<bapelParser::TupleExprContext>(0);
}


size_t bapelParser::AssignTermContext::getRuleIndex() const {
  return bapelParser::RuleAssignTerm;
}


antlrcpp::Any bapelParser::AssignTermContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitAssignTerm(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::AssignTermContext* bapelParser::assignTerm() {
  AssignTermContext *_localctx = _tracker.createInstance<AssignTermContext>(_ctx, getState());
  enterRule(_localctx, 94, bapelParser::RuleAssignTerm);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(572);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 54, _ctx)) {
    case 1: {
      setState(570);
      id();
      break;
    }

    case 2: {
      setState(571);
      tupleExpr();
      break;
    }

    default:
      break;
    }
    setState(574);
    match(bapelParser::LARROW);
    setState(575);
    expression();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ReturnTermContext ------------------------------------------------------------------

bapelParser::ReturnTermContext::ReturnTermContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::ReturnTermContext::RETURN() {
  return getToken(bapelParser::RETURN, 0);
}

bapelParser::ExpressionWithoutBlockContext* bapelParser::ReturnTermContext::expressionWithoutBlock() {
  return getRuleContext<bapelParser::ExpressionWithoutBlockContext>(0);
}


size_t bapelParser::ReturnTermContext::getRuleIndex() const {
  return bapelParser::RuleReturnTerm;
}


antlrcpp::Any bapelParser::ReturnTermContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitReturnTerm(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::ReturnTermContext* bapelParser::returnTerm() {
  ReturnTermContext *_localctx = _tracker.createInstance<ReturnTermContext>(_ctx, getState());
  enterRule(_localctx, 96, bapelParser::RuleReturnTerm);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(577);
    match(bapelParser::RETURN);
    setState(578);
    expressionWithoutBlock();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- IfTermContext ------------------------------------------------------------------

bapelParser::IfTermContext::IfTermContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::IfTermContext::IF() {
  return getToken(bapelParser::IF, 0);
}

bapelParser::ExpressionWithoutBlockContext* bapelParser::IfTermContext::expressionWithoutBlock() {
  return getRuleContext<bapelParser::ExpressionWithoutBlockContext>(0);
}

std::vector<bapelParser::BlockExprContext *> bapelParser::IfTermContext::blockExpr() {
  return getRuleContexts<bapelParser::BlockExprContext>();
}

bapelParser::BlockExprContext* bapelParser::IfTermContext::blockExpr(size_t i) {
  return getRuleContext<bapelParser::BlockExprContext>(i);
}

tree::TerminalNode* bapelParser::IfTermContext::ELSE() {
  return getToken(bapelParser::ELSE, 0);
}

bapelParser::IfTermContext* bapelParser::IfTermContext::ifTerm() {
  return getRuleContext<bapelParser::IfTermContext>(0);
}


size_t bapelParser::IfTermContext::getRuleIndex() const {
  return bapelParser::RuleIfTerm;
}


antlrcpp::Any bapelParser::IfTermContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitIfTerm(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::IfTermContext* bapelParser::ifTerm() {
  IfTermContext *_localctx = _tracker.createInstance<IfTermContext>(_ctx, getState());
  enterRule(_localctx, 98, bapelParser::RuleIfTerm);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(580);
    match(bapelParser::IF);
    setState(581);
    expressionWithoutBlock();
    setState(582);
    blockExpr();
    setState(588);
    _errHandler->sync(this);

    _la = _input->LA(1);
    if (_la == bapelParser::ELSE) {
      setState(583);
      match(bapelParser::ELSE);
      setState(586);
      _errHandler->sync(this);
      switch (_input->LA(1)) {
        case bapelParser::LBRACE: {
          setState(584);
          blockExpr();
          break;
        }

        case bapelParser::IF: {
          setState(585);
          ifTerm();
          break;
        }

      default:
        throw NoViableAltException(this);
      }
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ForTermContext ------------------------------------------------------------------

bapelParser::ForTermContext::ForTermContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::ForTermContext::FOR() {
  return getToken(bapelParser::FOR, 0);
}

bapelParser::ExpressionWithoutBlockContext* bapelParser::ForTermContext::expressionWithoutBlock() {
  return getRuleContext<bapelParser::ExpressionWithoutBlockContext>(0);
}

bapelParser::BlockExprContext* bapelParser::ForTermContext::blockExpr() {
  return getRuleContext<bapelParser::BlockExprContext>(0);
}


size_t bapelParser::ForTermContext::getRuleIndex() const {
  return bapelParser::RuleForTerm;
}


antlrcpp::Any bapelParser::ForTermContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitForTerm(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::ForTermContext* bapelParser::forTerm() {
  ForTermContext *_localctx = _tracker.createInstance<ForTermContext>(_ctx, getState());
  enterRule(_localctx, 100, bapelParser::RuleForTerm);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(590);
    match(bapelParser::FOR);
    setState(591);
    expressionWithoutBlock();
    setState(592);
    blockExpr();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- LambdaTermContext ------------------------------------------------------------------

bapelParser::LambdaTermContext::LambdaTermContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::LambdaTermContext::FN() {
  return getToken(bapelParser::FN, 0);
}

bapelParser::FunctionArgsContext* bapelParser::LambdaTermContext::functionArgs() {
  return getRuleContext<bapelParser::FunctionArgsContext>(0);
}

bapelParser::BlockExprContext* bapelParser::LambdaTermContext::blockExpr() {
  return getRuleContext<bapelParser::BlockExprContext>(0);
}

bapelParser::TypeAbstractionContext* bapelParser::LambdaTermContext::typeAbstraction() {
  return getRuleContext<bapelParser::TypeAbstractionContext>(0);
}


size_t bapelParser::LambdaTermContext::getRuleIndex() const {
  return bapelParser::RuleLambdaTerm;
}


antlrcpp::Any bapelParser::LambdaTermContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitLambdaTerm(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::LambdaTermContext* bapelParser::lambdaTerm() {
  LambdaTermContext *_localctx = _tracker.createInstance<LambdaTermContext>(_ctx, getState());
  enterRule(_localctx, 102, bapelParser::RuleLambdaTerm);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(594);
    match(bapelParser::FN);
    setState(596);
    _errHandler->sync(this);

    _la = _input->LA(1);
    if (_la == bapelParser::LBRACKET) {
      setState(595);
      typeAbstraction();
    }
    setState(598);
    functionArgs();
    setState(599);
    blockExpr();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- MatchTermContext ------------------------------------------------------------------

bapelParser::MatchTermContext::MatchTermContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::MatchTermContext::MATCH() {
  return getToken(bapelParser::MATCH, 0);
}

bapelParser::ExpressionContext* bapelParser::MatchTermContext::expression() {
  return getRuleContext<bapelParser::ExpressionContext>(0);
}

tree::TerminalNode* bapelParser::MatchTermContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

bapelParser::MatchArmsContext* bapelParser::MatchTermContext::matchArms() {
  return getRuleContext<bapelParser::MatchArmsContext>(0);
}

tree::TerminalNode* bapelParser::MatchTermContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}

tree::TerminalNode* bapelParser::MatchTermContext::COMMA() {
  return getToken(bapelParser::COMMA, 0);
}


size_t bapelParser::MatchTermContext::getRuleIndex() const {
  return bapelParser::RuleMatchTerm;
}


antlrcpp::Any bapelParser::MatchTermContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitMatchTerm(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::MatchTermContext* bapelParser::matchTerm() {
  MatchTermContext *_localctx = _tracker.createInstance<MatchTermContext>(_ctx, getState());
  enterRule(_localctx, 104, bapelParser::RuleMatchTerm);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(601);
    match(bapelParser::MATCH);
    setState(602);
    expression();
    setState(603);
    match(bapelParser::LBRACE);
    setState(604);
    matchArms();

    setState(606);
    _errHandler->sync(this);

    _la = _input->LA(1);
    if (_la == bapelParser::COMMA) {
      setState(605);
      match(bapelParser::COMMA);
    }
    setState(608);
    match(bapelParser::RBRACE);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- MatchArmsContext ------------------------------------------------------------------

bapelParser::MatchArmsContext::MatchArmsContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

std::vector<bapelParser::MatchArmContext *> bapelParser::MatchArmsContext::matchArm() {
  return getRuleContexts<bapelParser::MatchArmContext>();
}

bapelParser::MatchArmContext* bapelParser::MatchArmsContext::matchArm(size_t i) {
  return getRuleContext<bapelParser::MatchArmContext>(i);
}

std::vector<tree::TerminalNode *> bapelParser::MatchArmsContext::COMMA() {
  return getTokens(bapelParser::COMMA);
}

tree::TerminalNode* bapelParser::MatchArmsContext::COMMA(size_t i) {
  return getToken(bapelParser::COMMA, i);
}


size_t bapelParser::MatchArmsContext::getRuleIndex() const {
  return bapelParser::RuleMatchArms;
}


antlrcpp::Any bapelParser::MatchArmsContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitMatchArms(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::MatchArmsContext* bapelParser::matchArms() {
  MatchArmsContext *_localctx = _tracker.createInstance<MatchArmsContext>(_ctx, getState());
  enterRule(_localctx, 106, bapelParser::RuleMatchArms);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(610);
    matchArm();
    setState(615);
    _errHandler->sync(this);
    alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 59, _ctx);
    while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER) {
      if (alt == 1) {
        setState(611);
        match(bapelParser::COMMA);
        setState(612);
        matchArm(); 
      }
      setState(617);
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 59, _ctx);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- MatchArmContext ------------------------------------------------------------------

bapelParser::MatchArmContext::MatchArmContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::IdContext* bapelParser::MatchArmContext::id() {
  return getRuleContext<bapelParser::IdContext>(0);
}

tree::TerminalNode* bapelParser::MatchArmContext::IDENTIFIER() {
  return getToken(bapelParser::IDENTIFIER, 0);
}

tree::TerminalNode* bapelParser::MatchArmContext::FAT_ARROW() {
  return getToken(bapelParser::FAT_ARROW, 0);
}

bapelParser::ExpressionContext* bapelParser::MatchArmContext::expression() {
  return getRuleContext<bapelParser::ExpressionContext>(0);
}


size_t bapelParser::MatchArmContext::getRuleIndex() const {
  return bapelParser::RuleMatchArm;
}


antlrcpp::Any bapelParser::MatchArmContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitMatchArm(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::MatchArmContext* bapelParser::matchArm() {
  MatchArmContext *_localctx = _tracker.createInstance<MatchArmContext>(_ctx, getState());
  enterRule(_localctx, 108, bapelParser::RuleMatchArm);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(618);
    id();
    setState(619);
    match(bapelParser::IDENTIFIER);
    setState(620);
    match(bapelParser::FAT_ARROW);
    setState(621);
    expression();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- SetTermContext ------------------------------------------------------------------

bapelParser::SetTermContext::SetTermContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::SetTermContext::SET() {
  return getToken(bapelParser::SET, 0);
}

bapelParser::ExpressionContext* bapelParser::SetTermContext::expression() {
  return getRuleContext<bapelParser::ExpressionContext>(0);
}

tree::TerminalNode* bapelParser::SetTermContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

bapelParser::LabelValuesContext* bapelParser::SetTermContext::labelValues() {
  return getRuleContext<bapelParser::LabelValuesContext>(0);
}

tree::TerminalNode* bapelParser::SetTermContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}

tree::TerminalNode* bapelParser::SetTermContext::COMMA() {
  return getToken(bapelParser::COMMA, 0);
}


size_t bapelParser::SetTermContext::getRuleIndex() const {
  return bapelParser::RuleSetTerm;
}


antlrcpp::Any bapelParser::SetTermContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitSetTerm(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::SetTermContext* bapelParser::setTerm() {
  SetTermContext *_localctx = _tracker.createInstance<SetTermContext>(_ctx, getState());
  enterRule(_localctx, 110, bapelParser::RuleSetTerm);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(623);
    match(bapelParser::SET);
    setState(624);
    expression();
    setState(625);
    match(bapelParser::LBRACE);
    setState(626);
    labelValues();

    setState(628);
    _errHandler->sync(this);

    _la = _input->LA(1);
    if (_la == bapelParser::COMMA) {
      setState(627);
      match(bapelParser::COMMA);
    }
    setState(630);
    match(bapelParser::RBRACE);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- BlockExprContext ------------------------------------------------------------------

bapelParser::BlockExprContext::BlockExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::BlockExprContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

bapelParser::BlockStatementsContext* bapelParser::BlockExprContext::blockStatements() {
  return getRuleContext<bapelParser::BlockStatementsContext>(0);
}

tree::TerminalNode* bapelParser::BlockExprContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}


size_t bapelParser::BlockExprContext::getRuleIndex() const {
  return bapelParser::RuleBlockExpr;
}


antlrcpp::Any bapelParser::BlockExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitBlockExpr(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::BlockExprContext* bapelParser::blockExpr() {
  BlockExprContext *_localctx = _tracker.createInstance<BlockExprContext>(_ctx, getState());
  enterRule(_localctx, 112, bapelParser::RuleBlockExpr);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(632);
    match(bapelParser::LBRACE);
    setState(633);
    blockStatements();
    setState(634);
    match(bapelParser::RBRACE);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- BlockStatementsContext ------------------------------------------------------------------

bapelParser::BlockStatementsContext::BlockStatementsContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::StatementsContext* bapelParser::BlockStatementsContext::statements() {
  return getRuleContext<bapelParser::StatementsContext>(0);
}

bapelParser::ExpressionWithoutBlockContext* bapelParser::BlockStatementsContext::expressionWithoutBlock() {
  return getRuleContext<bapelParser::ExpressionWithoutBlockContext>(0);
}


size_t bapelParser::BlockStatementsContext::getRuleIndex() const {
  return bapelParser::RuleBlockStatements;
}


antlrcpp::Any bapelParser::BlockStatementsContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitBlockStatements(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::BlockStatementsContext* bapelParser::blockStatements() {
  BlockStatementsContext *_localctx = _tracker.createInstance<BlockStatementsContext>(_ctx, getState());
  enterRule(_localctx, 114, bapelParser::RuleBlockStatements);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(641);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 62, _ctx)) {
    case 1: {
      enterOuterAlt(_localctx, 1);
      setState(636);
      statements();
      setState(638);
      _errHandler->sync(this);

      _la = _input->LA(1);
      if ((((_la & ~ 0x3fULL) == 0) &&
        ((1ULL << _la) & ((1ULL << bapelParser::STRUCT)
        | (1ULL << bapelParser::VARIANT)
        | (1ULL << bapelParser::RETURN)
        | (1ULL << bapelParser::MINUS)
        | (1ULL << bapelParser::MUL)
        | (1ULL << bapelParser::NOT)
        | (1ULL << bapelParser::AMP)
        | (1ULL << bapelParser::LPAREN)
        | (1ULL << bapelParser::IDENTIFIER)
        | (1ULL << bapelParser::INT_LITERAL)
        | (1ULL << bapelParser::FLOAT_LITERAL)
        | (1ULL << bapelParser::RUNE_LITERAL)
        | (1ULL << bapelParser::STRING_LITERAL))) != 0)) {
        setState(637);
        expressionWithoutBlock();
      }
      break;
    }

    case 2: {
      enterOuterAlt(_localctx, 2);
      setState(640);
      expressionWithoutBlock();
      break;
    }

    default:
      break;
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- StatementsContext ------------------------------------------------------------------

bapelParser::StatementsContext::StatementsContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

std::vector<bapelParser::StatementContext *> bapelParser::StatementsContext::statement() {
  return getRuleContexts<bapelParser::StatementContext>();
}

bapelParser::StatementContext* bapelParser::StatementsContext::statement(size_t i) {
  return getRuleContext<bapelParser::StatementContext>(i);
}


size_t bapelParser::StatementsContext::getRuleIndex() const {
  return bapelParser::RuleStatements;
}


antlrcpp::Any bapelParser::StatementsContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitStatements(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::StatementsContext* bapelParser::statements() {
  StatementsContext *_localctx = _tracker.createInstance<StatementsContext>(_ctx, getState());
  enterRule(_localctx, 116, bapelParser::RuleStatements);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(644); 
    _errHandler->sync(this);
    alt = 1;
    do {
      switch (alt) {
        case 1: {
              setState(643);
              statement();
              break;
            }

      default:
        throw NoViableAltException(this);
      }
      setState(646); 
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 63, _ctx);
    } while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- StatementContext ------------------------------------------------------------------

bapelParser::StatementContext::StatementContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::LetStatementContext* bapelParser::StatementContext::letStatement() {
  return getRuleContext<bapelParser::LetStatementContext>(0);
}

bapelParser::ExpressionStatementContext* bapelParser::StatementContext::expressionStatement() {
  return getRuleContext<bapelParser::ExpressionStatementContext>(0);
}


size_t bapelParser::StatementContext::getRuleIndex() const {
  return bapelParser::RuleStatement;
}


antlrcpp::Any bapelParser::StatementContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitStatement(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::StatementContext* bapelParser::statement() {
  StatementContext *_localctx = _tracker.createInstance<StatementContext>(_ctx, getState());
  enterRule(_localctx, 118, bapelParser::RuleStatement);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(650);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::LET: {
        enterOuterAlt(_localctx, 1);
        setState(648);
        letStatement();
        break;
      }

      case bapelParser::FN:
      case bapelParser::STRUCT:
      case bapelParser::VARIANT:
      case bapelParser::MATCH:
      case bapelParser::SET:
      case bapelParser::RETURN:
      case bapelParser::IF:
      case bapelParser::FOR:
      case bapelParser::MINUS:
      case bapelParser::MUL:
      case bapelParser::NOT:
      case bapelParser::AMP:
      case bapelParser::LBRACE:
      case bapelParser::LPAREN:
      case bapelParser::IDENTIFIER:
      case bapelParser::INT_LITERAL:
      case bapelParser::FLOAT_LITERAL:
      case bapelParser::RUNE_LITERAL:
      case bapelParser::STRING_LITERAL: {
        enterOuterAlt(_localctx, 2);
        setState(649);
        expressionStatement();
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- LetStatementContext ------------------------------------------------------------------

bapelParser::LetStatementContext::LetStatementContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::LetStatementContext::LET() {
  return getToken(bapelParser::LET, 0);
}

bapelParser::IdContext* bapelParser::LetStatementContext::id() {
  return getRuleContext<bapelParser::IdContext>(0);
}

tree::TerminalNode* bapelParser::LetStatementContext::COLON() {
  return getToken(bapelParser::COLON, 0);
}

bapelParser::Type_Context* bapelParser::LetStatementContext::type_() {
  return getRuleContext<bapelParser::Type_Context>(0);
}

tree::TerminalNode* bapelParser::LetStatementContext::ASSIGN() {
  return getToken(bapelParser::ASSIGN, 0);
}

bapelParser::ExpressionContext* bapelParser::LetStatementContext::expression() {
  return getRuleContext<bapelParser::ExpressionContext>(0);
}

tree::TerminalNode* bapelParser::LetStatementContext::SEMI() {
  return getToken(bapelParser::SEMI, 0);
}


size_t bapelParser::LetStatementContext::getRuleIndex() const {
  return bapelParser::RuleLetStatement;
}


antlrcpp::Any bapelParser::LetStatementContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitLetStatement(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::LetStatementContext* bapelParser::letStatement() {
  LetStatementContext *_localctx = _tracker.createInstance<LetStatementContext>(_ctx, getState());
  enterRule(_localctx, 120, bapelParser::RuleLetStatement);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(666);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 65, _ctx)) {
    case 1: {
      enterOuterAlt(_localctx, 1);
      setState(652);
      match(bapelParser::LET);
      setState(653);
      id();
      setState(654);
      match(bapelParser::COLON);
      setState(655);
      type_();
      setState(656);
      match(bapelParser::ASSIGN);
      setState(657);
      expression();
      setState(658);
      match(bapelParser::SEMI);
      break;
    }

    case 2: {
      enterOuterAlt(_localctx, 2);
      setState(660);
      match(bapelParser::LET);
      setState(661);
      id();
      setState(662);
      match(bapelParser::ASSIGN);
      setState(663);
      expression();
      setState(664);
      match(bapelParser::SEMI);
      break;
    }

    default:
      break;
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ExpressionStatementContext ------------------------------------------------------------------

bapelParser::ExpressionStatementContext::ExpressionStatementContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::ExpressionWithoutBlockContext* bapelParser::ExpressionStatementContext::expressionWithoutBlock() {
  return getRuleContext<bapelParser::ExpressionWithoutBlockContext>(0);
}

tree::TerminalNode* bapelParser::ExpressionStatementContext::SEMI() {
  return getToken(bapelParser::SEMI, 0);
}

bapelParser::ExpressionWithBlockContext* bapelParser::ExpressionStatementContext::expressionWithBlock() {
  return getRuleContext<bapelParser::ExpressionWithBlockContext>(0);
}


size_t bapelParser::ExpressionStatementContext::getRuleIndex() const {
  return bapelParser::RuleExpressionStatement;
}


antlrcpp::Any bapelParser::ExpressionStatementContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitExpressionStatement(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::ExpressionStatementContext* bapelParser::expressionStatement() {
  ExpressionStatementContext *_localctx = _tracker.createInstance<ExpressionStatementContext>(_ctx, getState());
  enterRule(_localctx, 122, bapelParser::RuleExpressionStatement);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(675);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::STRUCT:
      case bapelParser::VARIANT:
      case bapelParser::RETURN:
      case bapelParser::MINUS:
      case bapelParser::MUL:
      case bapelParser::NOT:
      case bapelParser::AMP:
      case bapelParser::LPAREN:
      case bapelParser::IDENTIFIER:
      case bapelParser::INT_LITERAL:
      case bapelParser::FLOAT_LITERAL:
      case bapelParser::RUNE_LITERAL:
      case bapelParser::STRING_LITERAL: {
        enterOuterAlt(_localctx, 1);
        setState(668);
        expressionWithoutBlock();
        setState(669);
        match(bapelParser::SEMI);
        break;
      }

      case bapelParser::FN:
      case bapelParser::MATCH:
      case bapelParser::SET:
      case bapelParser::IF:
      case bapelParser::FOR:
      case bapelParser::LBRACE: {
        enterOuterAlt(_localctx, 2);
        setState(671);
        expressionWithBlock();
        setState(673);
        _errHandler->sync(this);

        _la = _input->LA(1);
        if (_la == bapelParser::SEMI) {
          setState(672);
          match(bapelParser::SEMI);
        }
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- OperatorExprContext ------------------------------------------------------------------

bapelParser::OperatorExprContext::OperatorExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::LogicalOrExprContext* bapelParser::OperatorExprContext::logicalOrExpr() {
  return getRuleContext<bapelParser::LogicalOrExprContext>(0);
}


size_t bapelParser::OperatorExprContext::getRuleIndex() const {
  return bapelParser::RuleOperatorExpr;
}


antlrcpp::Any bapelParser::OperatorExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitOperatorExpr(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::OperatorExprContext* bapelParser::operatorExpr() {
  OperatorExprContext *_localctx = _tracker.createInstance<OperatorExprContext>(_ctx, getState());
  enterRule(_localctx, 124, bapelParser::RuleOperatorExpr);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(677);
    logicalOrExpr(0);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- LogicalOrExprContext ------------------------------------------------------------------

bapelParser::LogicalOrExprContext::LogicalOrExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::LogicalAndExprContext* bapelParser::LogicalOrExprContext::logicalAndExpr() {
  return getRuleContext<bapelParser::LogicalAndExprContext>(0);
}

bapelParser::LogicalOrExprContext* bapelParser::LogicalOrExprContext::logicalOrExpr() {
  return getRuleContext<bapelParser::LogicalOrExprContext>(0);
}

tree::TerminalNode* bapelParser::LogicalOrExprContext::OR() {
  return getToken(bapelParser::OR, 0);
}


size_t bapelParser::LogicalOrExprContext::getRuleIndex() const {
  return bapelParser::RuleLogicalOrExpr;
}


antlrcpp::Any bapelParser::LogicalOrExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitLogicalOrExpr(this);
  else
    return visitor->visitChildren(this);
}


bapelParser::LogicalOrExprContext* bapelParser::logicalOrExpr() {
   return logicalOrExpr(0);
}

bapelParser::LogicalOrExprContext* bapelParser::logicalOrExpr(int precedence) {
  ParserRuleContext *parentContext = _ctx;
  size_t parentState = getState();
  bapelParser::LogicalOrExprContext *_localctx = _tracker.createInstance<LogicalOrExprContext>(_ctx, parentState);
  bapelParser::LogicalOrExprContext *previousContext = _localctx;
  (void)previousContext; // Silence compiler, in case the context is not used by generated code.
  size_t startState = 126;
  enterRecursionRule(_localctx, 126, bapelParser::RuleLogicalOrExpr, precedence);

    

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    unrollRecursionContexts(parentContext);
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(680);
    logicalAndExpr(0);
    _ctx->stop = _input->LT(-1);
    setState(687);
    _errHandler->sync(this);
    alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 68, _ctx);
    while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER) {
      if (alt == 1) {
        if (!_parseListeners.empty())
          triggerExitRuleEvent();
        previousContext = _localctx;
        _localctx = _tracker.createInstance<LogicalOrExprContext>(parentContext, parentState);
        pushNewRecursionContext(_localctx, startState, RuleLogicalOrExpr);
        setState(682);

        if (!(precpred(_ctx, 2))) throw FailedPredicateException(this, "precpred(_ctx, 2)");
        setState(683);
        match(bapelParser::OR);
        setState(684);
        logicalAndExpr(0); 
      }
      setState(689);
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 68, _ctx);
    }
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }
  return _localctx;
}

//----------------- LogicalAndExprContext ------------------------------------------------------------------

bapelParser::LogicalAndExprContext::LogicalAndExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::EqualityExprContext* bapelParser::LogicalAndExprContext::equalityExpr() {
  return getRuleContext<bapelParser::EqualityExprContext>(0);
}

bapelParser::LogicalAndExprContext* bapelParser::LogicalAndExprContext::logicalAndExpr() {
  return getRuleContext<bapelParser::LogicalAndExprContext>(0);
}

tree::TerminalNode* bapelParser::LogicalAndExprContext::AND() {
  return getToken(bapelParser::AND, 0);
}


size_t bapelParser::LogicalAndExprContext::getRuleIndex() const {
  return bapelParser::RuleLogicalAndExpr;
}


antlrcpp::Any bapelParser::LogicalAndExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitLogicalAndExpr(this);
  else
    return visitor->visitChildren(this);
}


bapelParser::LogicalAndExprContext* bapelParser::logicalAndExpr() {
   return logicalAndExpr(0);
}

bapelParser::LogicalAndExprContext* bapelParser::logicalAndExpr(int precedence) {
  ParserRuleContext *parentContext = _ctx;
  size_t parentState = getState();
  bapelParser::LogicalAndExprContext *_localctx = _tracker.createInstance<LogicalAndExprContext>(_ctx, parentState);
  bapelParser::LogicalAndExprContext *previousContext = _localctx;
  (void)previousContext; // Silence compiler, in case the context is not used by generated code.
  size_t startState = 128;
  enterRecursionRule(_localctx, 128, bapelParser::RuleLogicalAndExpr, precedence);

    

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    unrollRecursionContexts(parentContext);
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(691);
    equalityExpr(0);
    _ctx->stop = _input->LT(-1);
    setState(698);
    _errHandler->sync(this);
    alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 69, _ctx);
    while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER) {
      if (alt == 1) {
        if (!_parseListeners.empty())
          triggerExitRuleEvent();
        previousContext = _localctx;
        _localctx = _tracker.createInstance<LogicalAndExprContext>(parentContext, parentState);
        pushNewRecursionContext(_localctx, startState, RuleLogicalAndExpr);
        setState(693);

        if (!(precpred(_ctx, 2))) throw FailedPredicateException(this, "precpred(_ctx, 2)");
        setState(694);
        match(bapelParser::AND);
        setState(695);
        equalityExpr(0); 
      }
      setState(700);
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 69, _ctx);
    }
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }
  return _localctx;
}

//----------------- EqualityExprContext ------------------------------------------------------------------

bapelParser::EqualityExprContext::EqualityExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::ComparisonExprContext* bapelParser::EqualityExprContext::comparisonExpr() {
  return getRuleContext<bapelParser::ComparisonExprContext>(0);
}

bapelParser::EqualityExprContext* bapelParser::EqualityExprContext::equalityExpr() {
  return getRuleContext<bapelParser::EqualityExprContext>(0);
}

tree::TerminalNode* bapelParser::EqualityExprContext::NE() {
  return getToken(bapelParser::NE, 0);
}

tree::TerminalNode* bapelParser::EqualityExprContext::EQ() {
  return getToken(bapelParser::EQ, 0);
}

bapelParser::TypeApplicativeArgsContext* bapelParser::EqualityExprContext::typeApplicativeArgs() {
  return getRuleContext<bapelParser::TypeApplicativeArgsContext>(0);
}


size_t bapelParser::EqualityExprContext::getRuleIndex() const {
  return bapelParser::RuleEqualityExpr;
}


antlrcpp::Any bapelParser::EqualityExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitEqualityExpr(this);
  else
    return visitor->visitChildren(this);
}


bapelParser::EqualityExprContext* bapelParser::equalityExpr() {
   return equalityExpr(0);
}

bapelParser::EqualityExprContext* bapelParser::equalityExpr(int precedence) {
  ParserRuleContext *parentContext = _ctx;
  size_t parentState = getState();
  bapelParser::EqualityExprContext *_localctx = _tracker.createInstance<EqualityExprContext>(_ctx, parentState);
  bapelParser::EqualityExprContext *previousContext = _localctx;
  (void)previousContext; // Silence compiler, in case the context is not used by generated code.
  size_t startState = 130;
  enterRecursionRule(_localctx, 130, bapelParser::RuleEqualityExpr, precedence);

    size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    unrollRecursionContexts(parentContext);
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(702);
    comparisonExpr(0);
    _ctx->stop = _input->LT(-1);
    setState(712);
    _errHandler->sync(this);
    alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 71, _ctx);
    while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER) {
      if (alt == 1) {
        if (!_parseListeners.empty())
          triggerExitRuleEvent();
        previousContext = _localctx;
        _localctx = _tracker.createInstance<EqualityExprContext>(parentContext, parentState);
        pushNewRecursionContext(_localctx, startState, RuleEqualityExpr);
        setState(704);

        if (!(precpred(_ctx, 2))) throw FailedPredicateException(this, "precpred(_ctx, 2)");
        setState(705);
        _la = _input->LA(1);
        if (!(_la == bapelParser::NE

        || _la == bapelParser::EQ)) {
        _errHandler->recoverInline(this);
        }
        else {
          _errHandler->reportMatch(this);
          consume();
        }
        setState(707);
        _errHandler->sync(this);

        _la = _input->LA(1);
        if (_la == bapelParser::LBRACKET) {
          setState(706);
          typeApplicativeArgs();
        }
        setState(709);
        comparisonExpr(0); 
      }
      setState(714);
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 71, _ctx);
    }
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }
  return _localctx;
}

//----------------- ComparisonExprContext ------------------------------------------------------------------

bapelParser::ComparisonExprContext::ComparisonExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::AdditiveExprContext* bapelParser::ComparisonExprContext::additiveExpr() {
  return getRuleContext<bapelParser::AdditiveExprContext>(0);
}

bapelParser::ComparisonExprContext* bapelParser::ComparisonExprContext::comparisonExpr() {
  return getRuleContext<bapelParser::ComparisonExprContext>(0);
}

tree::TerminalNode* bapelParser::ComparisonExprContext::GT() {
  return getToken(bapelParser::GT, 0);
}

tree::TerminalNode* bapelParser::ComparisonExprContext::GE() {
  return getToken(bapelParser::GE, 0);
}

tree::TerminalNode* bapelParser::ComparisonExprContext::LT() {
  return getToken(bapelParser::LT, 0);
}

tree::TerminalNode* bapelParser::ComparisonExprContext::LE() {
  return getToken(bapelParser::LE, 0);
}

bapelParser::TypeApplicativeArgsContext* bapelParser::ComparisonExprContext::typeApplicativeArgs() {
  return getRuleContext<bapelParser::TypeApplicativeArgsContext>(0);
}


size_t bapelParser::ComparisonExprContext::getRuleIndex() const {
  return bapelParser::RuleComparisonExpr;
}


antlrcpp::Any bapelParser::ComparisonExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitComparisonExpr(this);
  else
    return visitor->visitChildren(this);
}


bapelParser::ComparisonExprContext* bapelParser::comparisonExpr() {
   return comparisonExpr(0);
}

bapelParser::ComparisonExprContext* bapelParser::comparisonExpr(int precedence) {
  ParserRuleContext *parentContext = _ctx;
  size_t parentState = getState();
  bapelParser::ComparisonExprContext *_localctx = _tracker.createInstance<ComparisonExprContext>(_ctx, parentState);
  bapelParser::ComparisonExprContext *previousContext = _localctx;
  (void)previousContext; // Silence compiler, in case the context is not used by generated code.
  size_t startState = 132;
  enterRecursionRule(_localctx, 132, bapelParser::RuleComparisonExpr, precedence);

    size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    unrollRecursionContexts(parentContext);
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(716);
    additiveExpr(0);
    _ctx->stop = _input->LT(-1);
    setState(726);
    _errHandler->sync(this);
    alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 73, _ctx);
    while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER) {
      if (alt == 1) {
        if (!_parseListeners.empty())
          triggerExitRuleEvent();
        previousContext = _localctx;
        _localctx = _tracker.createInstance<ComparisonExprContext>(parentContext, parentState);
        pushNewRecursionContext(_localctx, startState, RuleComparisonExpr);
        setState(718);

        if (!(precpred(_ctx, 2))) throw FailedPredicateException(this, "precpred(_ctx, 2)");
        setState(719);
        _la = _input->LA(1);
        if (!((((_la & ~ 0x3fULL) == 0) &&
          ((1ULL << _la) & ((1ULL << bapelParser::GE)
          | (1ULL << bapelParser::LE)
          | (1ULL << bapelParser::GT)
          | (1ULL << bapelParser::LT))) != 0))) {
        _errHandler->recoverInline(this);
        }
        else {
          _errHandler->reportMatch(this);
          consume();
        }
        setState(721);
        _errHandler->sync(this);

        _la = _input->LA(1);
        if (_la == bapelParser::LBRACKET) {
          setState(720);
          typeApplicativeArgs();
        }
        setState(723);
        additiveExpr(0); 
      }
      setState(728);
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 73, _ctx);
    }
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }
  return _localctx;
}

//----------------- AdditiveExprContext ------------------------------------------------------------------

bapelParser::AdditiveExprContext::AdditiveExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::MultiplicativeExprContext* bapelParser::AdditiveExprContext::multiplicativeExpr() {
  return getRuleContext<bapelParser::MultiplicativeExprContext>(0);
}

bapelParser::AdditiveExprContext* bapelParser::AdditiveExprContext::additiveExpr() {
  return getRuleContext<bapelParser::AdditiveExprContext>(0);
}

tree::TerminalNode* bapelParser::AdditiveExprContext::PLUS() {
  return getToken(bapelParser::PLUS, 0);
}

tree::TerminalNode* bapelParser::AdditiveExprContext::MINUS() {
  return getToken(bapelParser::MINUS, 0);
}

bapelParser::TypeApplicativeArgsContext* bapelParser::AdditiveExprContext::typeApplicativeArgs() {
  return getRuleContext<bapelParser::TypeApplicativeArgsContext>(0);
}


size_t bapelParser::AdditiveExprContext::getRuleIndex() const {
  return bapelParser::RuleAdditiveExpr;
}


antlrcpp::Any bapelParser::AdditiveExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitAdditiveExpr(this);
  else
    return visitor->visitChildren(this);
}


bapelParser::AdditiveExprContext* bapelParser::additiveExpr() {
   return additiveExpr(0);
}

bapelParser::AdditiveExprContext* bapelParser::additiveExpr(int precedence) {
  ParserRuleContext *parentContext = _ctx;
  size_t parentState = getState();
  bapelParser::AdditiveExprContext *_localctx = _tracker.createInstance<AdditiveExprContext>(_ctx, parentState);
  bapelParser::AdditiveExprContext *previousContext = _localctx;
  (void)previousContext; // Silence compiler, in case the context is not used by generated code.
  size_t startState = 134;
  enterRecursionRule(_localctx, 134, bapelParser::RuleAdditiveExpr, precedence);

    size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    unrollRecursionContexts(parentContext);
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(730);
    multiplicativeExpr(0);
    _ctx->stop = _input->LT(-1);
    setState(740);
    _errHandler->sync(this);
    alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 75, _ctx);
    while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER) {
      if (alt == 1) {
        if (!_parseListeners.empty())
          triggerExitRuleEvent();
        previousContext = _localctx;
        _localctx = _tracker.createInstance<AdditiveExprContext>(parentContext, parentState);
        pushNewRecursionContext(_localctx, startState, RuleAdditiveExpr);
        setState(732);

        if (!(precpred(_ctx, 2))) throw FailedPredicateException(this, "precpred(_ctx, 2)");
        setState(733);
        _la = _input->LA(1);
        if (!(_la == bapelParser::PLUS

        || _la == bapelParser::MINUS)) {
        _errHandler->recoverInline(this);
        }
        else {
          _errHandler->reportMatch(this);
          consume();
        }
        setState(735);
        _errHandler->sync(this);

        _la = _input->LA(1);
        if (_la == bapelParser::LBRACKET) {
          setState(734);
          typeApplicativeArgs();
        }
        setState(737);
        multiplicativeExpr(0); 
      }
      setState(742);
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 75, _ctx);
    }
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }
  return _localctx;
}

//----------------- MultiplicativeExprContext ------------------------------------------------------------------

bapelParser::MultiplicativeExprContext::MultiplicativeExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::UnaryExprContext* bapelParser::MultiplicativeExprContext::unaryExpr() {
  return getRuleContext<bapelParser::UnaryExprContext>(0);
}

bapelParser::MultiplicativeExprContext* bapelParser::MultiplicativeExprContext::multiplicativeExpr() {
  return getRuleContext<bapelParser::MultiplicativeExprContext>(0);
}

tree::TerminalNode* bapelParser::MultiplicativeExprContext::MUL() {
  return getToken(bapelParser::MUL, 0);
}

tree::TerminalNode* bapelParser::MultiplicativeExprContext::DIV() {
  return getToken(bapelParser::DIV, 0);
}

bapelParser::TypeApplicativeArgsContext* bapelParser::MultiplicativeExprContext::typeApplicativeArgs() {
  return getRuleContext<bapelParser::TypeApplicativeArgsContext>(0);
}


size_t bapelParser::MultiplicativeExprContext::getRuleIndex() const {
  return bapelParser::RuleMultiplicativeExpr;
}


antlrcpp::Any bapelParser::MultiplicativeExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitMultiplicativeExpr(this);
  else
    return visitor->visitChildren(this);
}


bapelParser::MultiplicativeExprContext* bapelParser::multiplicativeExpr() {
   return multiplicativeExpr(0);
}

bapelParser::MultiplicativeExprContext* bapelParser::multiplicativeExpr(int precedence) {
  ParserRuleContext *parentContext = _ctx;
  size_t parentState = getState();
  bapelParser::MultiplicativeExprContext *_localctx = _tracker.createInstance<MultiplicativeExprContext>(_ctx, parentState);
  bapelParser::MultiplicativeExprContext *previousContext = _localctx;
  (void)previousContext; // Silence compiler, in case the context is not used by generated code.
  size_t startState = 136;
  enterRecursionRule(_localctx, 136, bapelParser::RuleMultiplicativeExpr, precedence);

    size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    unrollRecursionContexts(parentContext);
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(744);
    unaryExpr();
    _ctx->stop = _input->LT(-1);
    setState(754);
    _errHandler->sync(this);
    alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 77, _ctx);
    while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER) {
      if (alt == 1) {
        if (!_parseListeners.empty())
          triggerExitRuleEvent();
        previousContext = _localctx;
        _localctx = _tracker.createInstance<MultiplicativeExprContext>(parentContext, parentState);
        pushNewRecursionContext(_localctx, startState, RuleMultiplicativeExpr);
        setState(746);

        if (!(precpred(_ctx, 2))) throw FailedPredicateException(this, "precpred(_ctx, 2)");
        setState(747);
        _la = _input->LA(1);
        if (!(_la == bapelParser::MUL

        || _la == bapelParser::DIV)) {
        _errHandler->recoverInline(this);
        }
        else {
          _errHandler->reportMatch(this);
          consume();
        }
        setState(749);
        _errHandler->sync(this);

        _la = _input->LA(1);
        if (_la == bapelParser::LBRACKET) {
          setState(748);
          typeApplicativeArgs();
        }
        setState(751);
        unaryExpr(); 
      }
      setState(756);
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 77, _ctx);
    }
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }
  return _localctx;
}

//----------------- UnaryExprContext ------------------------------------------------------------------

bapelParser::UnaryExprContext::UnaryExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::UnaryExprContext* bapelParser::UnaryExprContext::unaryExpr() {
  return getRuleContext<bapelParser::UnaryExprContext>(0);
}

tree::TerminalNode* bapelParser::UnaryExprContext::NOT() {
  return getToken(bapelParser::NOT, 0);
}

tree::TerminalNode* bapelParser::UnaryExprContext::MINUS() {
  return getToken(bapelParser::MINUS, 0);
}

bapelParser::TypeApplicativeArgsContext* bapelParser::UnaryExprContext::typeApplicativeArgs() {
  return getRuleContext<bapelParser::TypeApplicativeArgsContext>(0);
}

bapelParser::ApplicativeExprContext* bapelParser::UnaryExprContext::applicativeExpr() {
  return getRuleContext<bapelParser::ApplicativeExprContext>(0);
}


size_t bapelParser::UnaryExprContext::getRuleIndex() const {
  return bapelParser::RuleUnaryExpr;
}


antlrcpp::Any bapelParser::UnaryExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitUnaryExpr(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::UnaryExprContext* bapelParser::unaryExpr() {
  UnaryExprContext *_localctx = _tracker.createInstance<UnaryExprContext>(_ctx, getState());
  enterRule(_localctx, 138, bapelParser::RuleUnaryExpr);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(763);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::MINUS:
      case bapelParser::NOT: {
        enterOuterAlt(_localctx, 1);
        setState(757);
        _la = _input->LA(1);
        if (!(_la == bapelParser::MINUS

        || _la == bapelParser::NOT)) {
        _errHandler->recoverInline(this);
        }
        else {
          _errHandler->reportMatch(this);
          consume();
        }
        setState(759);
        _errHandler->sync(this);

        _la = _input->LA(1);
        if (_la == bapelParser::LBRACKET) {
          setState(758);
          typeApplicativeArgs();
        }
        setState(761);
        unaryExpr();
        break;
      }

      case bapelParser::STRUCT:
      case bapelParser::VARIANT:
      case bapelParser::MUL:
      case bapelParser::AMP:
      case bapelParser::LPAREN:
      case bapelParser::IDENTIFIER:
      case bapelParser::INT_LITERAL:
      case bapelParser::FLOAT_LITERAL:
      case bapelParser::RUNE_LITERAL:
      case bapelParser::STRING_LITERAL: {
        enterOuterAlt(_localctx, 2);
        setState(762);
        applicativeExpr(0);
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ApplicativeExprContext ------------------------------------------------------------------

bapelParser::ApplicativeExprContext::ApplicativeExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::TypeApplicativeExprContext* bapelParser::ApplicativeExprContext::typeApplicativeExpr() {
  return getRuleContext<bapelParser::TypeApplicativeExprContext>(0);
}

bapelParser::ApplicativeExprContext* bapelParser::ApplicativeExprContext::applicativeExpr() {
  return getRuleContext<bapelParser::ApplicativeExprContext>(0);
}

bapelParser::BasePrimaryExprContext* bapelParser::ApplicativeExprContext::basePrimaryExpr() {
  return getRuleContext<bapelParser::BasePrimaryExprContext>(0);
}


size_t bapelParser::ApplicativeExprContext::getRuleIndex() const {
  return bapelParser::RuleApplicativeExpr;
}


antlrcpp::Any bapelParser::ApplicativeExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitApplicativeExpr(this);
  else
    return visitor->visitChildren(this);
}


bapelParser::ApplicativeExprContext* bapelParser::applicativeExpr() {
   return applicativeExpr(0);
}

bapelParser::ApplicativeExprContext* bapelParser::applicativeExpr(int precedence) {
  ParserRuleContext *parentContext = _ctx;
  size_t parentState = getState();
  bapelParser::ApplicativeExprContext *_localctx = _tracker.createInstance<ApplicativeExprContext>(_ctx, parentState);
  bapelParser::ApplicativeExprContext *previousContext = _localctx;
  (void)previousContext; // Silence compiler, in case the context is not used by generated code.
  size_t startState = 140;
  enterRecursionRule(_localctx, 140, bapelParser::RuleApplicativeExpr, precedence);

    

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    unrollRecursionContexts(parentContext);
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(766);
    typeApplicativeExpr();
    _ctx->stop = _input->LT(-1);
    setState(772);
    _errHandler->sync(this);
    alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 80, _ctx);
    while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER) {
      if (alt == 1) {
        if (!_parseListeners.empty())
          triggerExitRuleEvent();
        previousContext = _localctx;
        _localctx = _tracker.createInstance<ApplicativeExprContext>(parentContext, parentState);
        pushNewRecursionContext(_localctx, startState, RuleApplicativeExpr);
        setState(768);

        if (!(precpred(_ctx, 2))) throw FailedPredicateException(this, "precpred(_ctx, 2)");
        setState(769);
        basePrimaryExpr(); 
      }
      setState(774);
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 80, _ctx);
    }
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }
  return _localctx;
}

//----------------- TypeApplicativeExprContext ------------------------------------------------------------------

bapelParser::TypeApplicativeExprContext::TypeApplicativeExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::PrimaryExprContext* bapelParser::TypeApplicativeExprContext::primaryExpr() {
  return getRuleContext<bapelParser::PrimaryExprContext>(0);
}

bapelParser::TypeApplicativeArgsContext* bapelParser::TypeApplicativeExprContext::typeApplicativeArgs() {
  return getRuleContext<bapelParser::TypeApplicativeArgsContext>(0);
}


size_t bapelParser::TypeApplicativeExprContext::getRuleIndex() const {
  return bapelParser::RuleTypeApplicativeExpr;
}


antlrcpp::Any bapelParser::TypeApplicativeExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTypeApplicativeExpr(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TypeApplicativeExprContext* bapelParser::typeApplicativeExpr() {
  TypeApplicativeExprContext *_localctx = _tracker.createInstance<TypeApplicativeExprContext>(_ctx, getState());
  enterRule(_localctx, 142, bapelParser::RuleTypeApplicativeExpr);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(779);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 81, _ctx)) {
    case 1: {
      enterOuterAlt(_localctx, 1);
      setState(775);
      primaryExpr();
      setState(776);
      typeApplicativeArgs();
      break;
    }

    case 2: {
      enterOuterAlt(_localctx, 2);
      setState(778);
      primaryExpr();
      break;
    }

    default:
      break;
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- TypeApplicativeArgsContext ------------------------------------------------------------------

bapelParser::TypeApplicativeArgsContext::TypeApplicativeArgsContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::TypeApplicativeArgsContext::LBRACKET() {
  return getToken(bapelParser::LBRACKET, 0);
}

bapelParser::TupleTypeArgsContext* bapelParser::TypeApplicativeArgsContext::tupleTypeArgs() {
  return getRuleContext<bapelParser::TupleTypeArgsContext>(0);
}

tree::TerminalNode* bapelParser::TypeApplicativeArgsContext::RBRACKET() {
  return getToken(bapelParser::RBRACKET, 0);
}

bapelParser::Type_Context* bapelParser::TypeApplicativeArgsContext::type_() {
  return getRuleContext<bapelParser::Type_Context>(0);
}


size_t bapelParser::TypeApplicativeArgsContext::getRuleIndex() const {
  return bapelParser::RuleTypeApplicativeArgs;
}


antlrcpp::Any bapelParser::TypeApplicativeArgsContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTypeApplicativeArgs(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TypeApplicativeArgsContext* bapelParser::typeApplicativeArgs() {
  TypeApplicativeArgsContext *_localctx = _tracker.createInstance<TypeApplicativeArgsContext>(_ctx, getState());
  enterRule(_localctx, 144, bapelParser::RuleTypeApplicativeArgs);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(789);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 82, _ctx)) {
    case 1: {
      enterOuterAlt(_localctx, 1);
      setState(781);
      match(bapelParser::LBRACKET);
      setState(782);
      tupleTypeArgs();
      setState(783);
      match(bapelParser::RBRACKET);
      break;
    }

    case 2: {
      enterOuterAlt(_localctx, 2);
      setState(785);
      match(bapelParser::LBRACKET);
      setState(786);
      type_();
      setState(787);
      match(bapelParser::RBRACKET);
      break;
    }

    default:
      break;
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- BasePrimaryExprContext ------------------------------------------------------------------

bapelParser::BasePrimaryExprContext::BasePrimaryExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::BasePrimaryExprContext::AMP() {
  return getToken(bapelParser::AMP, 0);
}

bapelParser::ProjectionExprContext* bapelParser::BasePrimaryExprContext::projectionExpr() {
  return getRuleContext<bapelParser::ProjectionExprContext>(0);
}

tree::TerminalNode* bapelParser::BasePrimaryExprContext::INT_LITERAL() {
  return getToken(bapelParser::INT_LITERAL, 0);
}

tree::TerminalNode* bapelParser::BasePrimaryExprContext::FLOAT_LITERAL() {
  return getToken(bapelParser::FLOAT_LITERAL, 0);
}


size_t bapelParser::BasePrimaryExprContext::getRuleIndex() const {
  return bapelParser::RuleBasePrimaryExpr;
}


antlrcpp::Any bapelParser::BasePrimaryExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitBasePrimaryExpr(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::BasePrimaryExprContext* bapelParser::basePrimaryExpr() {
  BasePrimaryExprContext *_localctx = _tracker.createInstance<BasePrimaryExprContext>(_ctx, getState());
  enterRule(_localctx, 146, bapelParser::RuleBasePrimaryExpr);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(796);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::AMP: {
        enterOuterAlt(_localctx, 1);
        setState(791);
        match(bapelParser::AMP);
        setState(792);
        projectionExpr(0);
        break;
      }

      case bapelParser::STRUCT:
      case bapelParser::VARIANT:
      case bapelParser::LPAREN:
      case bapelParser::IDENTIFIER:
      case bapelParser::RUNE_LITERAL:
      case bapelParser::STRING_LITERAL: {
        enterOuterAlt(_localctx, 2);
        setState(793);
        projectionExpr(0);
        break;
      }

      case bapelParser::INT_LITERAL: {
        enterOuterAlt(_localctx, 3);
        setState(794);
        match(bapelParser::INT_LITERAL);
        break;
      }

      case bapelParser::FLOAT_LITERAL: {
        enterOuterAlt(_localctx, 4);
        setState(795);
        match(bapelParser::FLOAT_LITERAL);
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- PrimaryExprContext ------------------------------------------------------------------

bapelParser::PrimaryExprContext::PrimaryExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::PrimaryExprContext::MUL() {
  return getToken(bapelParser::MUL, 0);
}

bapelParser::PrimaryExprContext* bapelParser::PrimaryExprContext::primaryExpr() {
  return getRuleContext<bapelParser::PrimaryExprContext>(0);
}

bapelParser::BasePrimaryExprContext* bapelParser::PrimaryExprContext::basePrimaryExpr() {
  return getRuleContext<bapelParser::BasePrimaryExprContext>(0);
}


size_t bapelParser::PrimaryExprContext::getRuleIndex() const {
  return bapelParser::RulePrimaryExpr;
}


antlrcpp::Any bapelParser::PrimaryExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitPrimaryExpr(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::PrimaryExprContext* bapelParser::primaryExpr() {
  PrimaryExprContext *_localctx = _tracker.createInstance<PrimaryExprContext>(_ctx, getState());
  enterRule(_localctx, 148, bapelParser::RulePrimaryExpr);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(801);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::MUL: {
        enterOuterAlt(_localctx, 1);
        setState(798);
        match(bapelParser::MUL);
        setState(799);
        primaryExpr();
        break;
      }

      case bapelParser::STRUCT:
      case bapelParser::VARIANT:
      case bapelParser::AMP:
      case bapelParser::LPAREN:
      case bapelParser::IDENTIFIER:
      case bapelParser::INT_LITERAL:
      case bapelParser::FLOAT_LITERAL:
      case bapelParser::RUNE_LITERAL:
      case bapelParser::STRING_LITERAL: {
        enterOuterAlt(_localctx, 2);
        setState(800);
        basePrimaryExpr();
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- ProjectionExprContext ------------------------------------------------------------------

bapelParser::ProjectionExprContext::ProjectionExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::DerefExprContext* bapelParser::ProjectionExprContext::derefExpr() {
  return getRuleContext<bapelParser::DerefExprContext>(0);
}

bapelParser::ProjectionExprContext* bapelParser::ProjectionExprContext::projectionExpr() {
  return getRuleContext<bapelParser::ProjectionExprContext>(0);
}

tree::TerminalNode* bapelParser::ProjectionExprContext::DOT() {
  return getToken(bapelParser::DOT, 0);
}

tree::TerminalNode* bapelParser::ProjectionExprContext::INT_LITERAL() {
  return getToken(bapelParser::INT_LITERAL, 0);
}

tree::TerminalNode* bapelParser::ProjectionExprContext::IDENTIFIER() {
  return getToken(bapelParser::IDENTIFIER, 0);
}


size_t bapelParser::ProjectionExprContext::getRuleIndex() const {
  return bapelParser::RuleProjectionExpr;
}


antlrcpp::Any bapelParser::ProjectionExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitProjectionExpr(this);
  else
    return visitor->visitChildren(this);
}


bapelParser::ProjectionExprContext* bapelParser::projectionExpr() {
   return projectionExpr(0);
}

bapelParser::ProjectionExprContext* bapelParser::projectionExpr(int precedence) {
  ParserRuleContext *parentContext = _ctx;
  size_t parentState = getState();
  bapelParser::ProjectionExprContext *_localctx = _tracker.createInstance<ProjectionExprContext>(_ctx, parentState);
  bapelParser::ProjectionExprContext *previousContext = _localctx;
  (void)previousContext; // Silence compiler, in case the context is not used by generated code.
  size_t startState = 150;
  enterRecursionRule(_localctx, 150, bapelParser::RuleProjectionExpr, precedence);

    

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    unrollRecursionContexts(parentContext);
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(804);
    derefExpr();
    _ctx->stop = _input->LT(-1);
    setState(814);
    _errHandler->sync(this);
    alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 86, _ctx);
    while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER) {
      if (alt == 1) {
        if (!_parseListeners.empty())
          triggerExitRuleEvent();
        previousContext = _localctx;
        setState(812);
        _errHandler->sync(this);
        switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 85, _ctx)) {
        case 1: {
          _localctx = _tracker.createInstance<ProjectionExprContext>(parentContext, parentState);
          pushNewRecursionContext(_localctx, startState, RuleProjectionExpr);
          setState(806);

          if (!(precpred(_ctx, 3))) throw FailedPredicateException(this, "precpred(_ctx, 3)");
          setState(807);
          match(bapelParser::DOT);
          setState(808);
          match(bapelParser::INT_LITERAL);
          break;
        }

        case 2: {
          _localctx = _tracker.createInstance<ProjectionExprContext>(parentContext, parentState);
          pushNewRecursionContext(_localctx, startState, RuleProjectionExpr);
          setState(809);

          if (!(precpred(_ctx, 2))) throw FailedPredicateException(this, "precpred(_ctx, 2)");
          setState(810);
          match(bapelParser::DOT);
          setState(811);
          match(bapelParser::IDENTIFIER);
          break;
        }

        default:
          break;
        } 
      }
      setState(816);
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 86, _ctx);
    }
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }
  return _localctx;
}

//----------------- DerefExprContext ------------------------------------------------------------------

bapelParser::DerefExprContext::DerefExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::InjectionExprContext* bapelParser::DerefExprContext::injectionExpr() {
  return getRuleContext<bapelParser::InjectionExprContext>(0);
}

tree::TerminalNode* bapelParser::DerefExprContext::RUNE_LITERAL() {
  return getToken(bapelParser::RUNE_LITERAL, 0);
}

tree::TerminalNode* bapelParser::DerefExprContext::STRING_LITERAL() {
  return getToken(bapelParser::STRING_LITERAL, 0);
}

bapelParser::StructExprContext* bapelParser::DerefExprContext::structExpr() {
  return getRuleContext<bapelParser::StructExprContext>(0);
}

bapelParser::TupleExprContext* bapelParser::DerefExprContext::tupleExpr() {
  return getRuleContext<bapelParser::TupleExprContext>(0);
}

bapelParser::VarExprContext* bapelParser::DerefExprContext::varExpr() {
  return getRuleContext<bapelParser::VarExprContext>(0);
}

tree::TerminalNode* bapelParser::DerefExprContext::LPAREN() {
  return getToken(bapelParser::LPAREN, 0);
}

bapelParser::ExpressionContext* bapelParser::DerefExprContext::expression() {
  return getRuleContext<bapelParser::ExpressionContext>(0);
}

tree::TerminalNode* bapelParser::DerefExprContext::RPAREN() {
  return getToken(bapelParser::RPAREN, 0);
}


size_t bapelParser::DerefExprContext::getRuleIndex() const {
  return bapelParser::RuleDerefExpr;
}


antlrcpp::Any bapelParser::DerefExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitDerefExpr(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::DerefExprContext* bapelParser::derefExpr() {
  DerefExprContext *_localctx = _tracker.createInstance<DerefExprContext>(_ctx, getState());
  enterRule(_localctx, 152, bapelParser::RuleDerefExpr);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(827);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 87, _ctx)) {
    case 1: {
      enterOuterAlt(_localctx, 1);
      setState(817);
      injectionExpr();
      break;
    }

    case 2: {
      enterOuterAlt(_localctx, 2);
      setState(818);
      match(bapelParser::RUNE_LITERAL);
      break;
    }

    case 3: {
      enterOuterAlt(_localctx, 3);
      setState(819);
      match(bapelParser::STRING_LITERAL);
      break;
    }

    case 4: {
      enterOuterAlt(_localctx, 4);
      setState(820);
      structExpr();
      break;
    }

    case 5: {
      enterOuterAlt(_localctx, 5);
      setState(821);
      tupleExpr();
      break;
    }

    case 6: {
      enterOuterAlt(_localctx, 6);
      setState(822);
      varExpr();
      break;
    }

    case 7: {
      enterOuterAlt(_localctx, 7);
      setState(823);
      match(bapelParser::LPAREN);
      setState(824);
      expression();
      setState(825);
      match(bapelParser::RPAREN);
      break;
    }

    default:
      break;
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- InjectionExprContext ------------------------------------------------------------------

bapelParser::InjectionExprContext::InjectionExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::InjectionExprContext::VARIANT() {
  return getToken(bapelParser::VARIANT, 0);
}

tree::TerminalNode* bapelParser::InjectionExprContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

bapelParser::Type_Context* bapelParser::InjectionExprContext::type_() {
  return getRuleContext<bapelParser::Type_Context>(0);
}

bapelParser::LabelValueContext* bapelParser::InjectionExprContext::labelValue() {
  return getRuleContext<bapelParser::LabelValueContext>(0);
}

tree::TerminalNode* bapelParser::InjectionExprContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}


size_t bapelParser::InjectionExprContext::getRuleIndex() const {
  return bapelParser::RuleInjectionExpr;
}


antlrcpp::Any bapelParser::InjectionExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitInjectionExpr(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::InjectionExprContext* bapelParser::injectionExpr() {
  InjectionExprContext *_localctx = _tracker.createInstance<InjectionExprContext>(_ctx, getState());
  enterRule(_localctx, 154, bapelParser::RuleInjectionExpr);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(829);
    match(bapelParser::VARIANT);
    setState(830);
    match(bapelParser::LBRACE);
    setState(831);
    type_();
    setState(832);
    labelValue();
    setState(833);
    match(bapelParser::RBRACE);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- StructExprContext ------------------------------------------------------------------

bapelParser::StructExprContext::StructExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::StructExprContext::STRUCT() {
  return getToken(bapelParser::STRUCT, 0);
}

tree::TerminalNode* bapelParser::StructExprContext::LBRACE() {
  return getToken(bapelParser::LBRACE, 0);
}

tree::TerminalNode* bapelParser::StructExprContext::RBRACE() {
  return getToken(bapelParser::RBRACE, 0);
}

bapelParser::LabelValuesContext* bapelParser::StructExprContext::labelValues() {
  return getRuleContext<bapelParser::LabelValuesContext>(0);
}

tree::TerminalNode* bapelParser::StructExprContext::COMMA() {
  return getToken(bapelParser::COMMA, 0);
}


size_t bapelParser::StructExprContext::getRuleIndex() const {
  return bapelParser::RuleStructExpr;
}


antlrcpp::Any bapelParser::StructExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitStructExpr(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::StructExprContext* bapelParser::structExpr() {
  StructExprContext *_localctx = _tracker.createInstance<StructExprContext>(_ctx, getState());
  enterRule(_localctx, 156, bapelParser::RuleStructExpr);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(835);
    match(bapelParser::STRUCT);
    setState(836);
    match(bapelParser::LBRACE);
    setState(841);
    _errHandler->sync(this);

    _la = _input->LA(1);
    if ((((_la & ~ 0x3fULL) == 0) &&
      ((1ULL << _la) & ((1ULL << bapelParser::LPAREN)
      | (1ULL << bapelParser::IDENTIFIER)
      | (1ULL << bapelParser::INT_LITERAL))) != 0)) {
      setState(837);
      labelValues();

      setState(839);
      _errHandler->sync(this);

      _la = _input->LA(1);
      if (_la == bapelParser::COMMA) {
        setState(838);
        match(bapelParser::COMMA);
      }
    }
    setState(843);
    match(bapelParser::RBRACE);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- LabelValuesContext ------------------------------------------------------------------

bapelParser::LabelValuesContext::LabelValuesContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

std::vector<bapelParser::LabelValueContext *> bapelParser::LabelValuesContext::labelValue() {
  return getRuleContexts<bapelParser::LabelValueContext>();
}

bapelParser::LabelValueContext* bapelParser::LabelValuesContext::labelValue(size_t i) {
  return getRuleContext<bapelParser::LabelValueContext>(i);
}

std::vector<tree::TerminalNode *> bapelParser::LabelValuesContext::COMMA() {
  return getTokens(bapelParser::COMMA);
}

tree::TerminalNode* bapelParser::LabelValuesContext::COMMA(size_t i) {
  return getToken(bapelParser::COMMA, i);
}


size_t bapelParser::LabelValuesContext::getRuleIndex() const {
  return bapelParser::RuleLabelValues;
}


antlrcpp::Any bapelParser::LabelValuesContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitLabelValues(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::LabelValuesContext* bapelParser::labelValues() {
  LabelValuesContext *_localctx = _tracker.createInstance<LabelValuesContext>(_ctx, getState());
  enterRule(_localctx, 158, bapelParser::RuleLabelValues);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(845);
    labelValue();
    setState(850);
    _errHandler->sync(this);
    alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 90, _ctx);
    while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER) {
      if (alt == 1) {
        setState(846);
        match(bapelParser::COMMA);
        setState(847);
        labelValue(); 
      }
      setState(852);
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 90, _ctx);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- LabelValueContext ------------------------------------------------------------------

bapelParser::LabelValueContext::LabelValueContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::IdContext* bapelParser::LabelValueContext::id() {
  return getRuleContext<bapelParser::IdContext>(0);
}

tree::TerminalNode* bapelParser::LabelValueContext::ASSIGN() {
  return getToken(bapelParser::ASSIGN, 0);
}

bapelParser::ExpressionContext* bapelParser::LabelValueContext::expression() {
  return getRuleContext<bapelParser::ExpressionContext>(0);
}

tree::TerminalNode* bapelParser::LabelValueContext::INT_LITERAL() {
  return getToken(bapelParser::INT_LITERAL, 0);
}


size_t bapelParser::LabelValueContext::getRuleIndex() const {
  return bapelParser::RuleLabelValue;
}


antlrcpp::Any bapelParser::LabelValueContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitLabelValue(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::LabelValueContext* bapelParser::labelValue() {
  LabelValueContext *_localctx = _tracker.createInstance<LabelValueContext>(_ctx, getState());
  enterRule(_localctx, 160, bapelParser::RuleLabelValue);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(860);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::LPAREN:
      case bapelParser::IDENTIFIER: {
        enterOuterAlt(_localctx, 1);
        setState(853);
        id();
        setState(854);
        match(bapelParser::ASSIGN);
        setState(855);
        expression();
        break;
      }

      case bapelParser::INT_LITERAL: {
        enterOuterAlt(_localctx, 2);
        setState(857);
        match(bapelParser::INT_LITERAL);
        setState(858);
        match(bapelParser::ASSIGN);
        setState(859);
        expression();
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- TupleExprContext ------------------------------------------------------------------

bapelParser::TupleExprContext::TupleExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

tree::TerminalNode* bapelParser::TupleExprContext::LPAREN() {
  return getToken(bapelParser::LPAREN, 0);
}

tree::TerminalNode* bapelParser::TupleExprContext::RPAREN() {
  return getToken(bapelParser::RPAREN, 0);
}

bapelParser::TupleExprArgsContext* bapelParser::TupleExprContext::tupleExprArgs() {
  return getRuleContext<bapelParser::TupleExprArgsContext>(0);
}


size_t bapelParser::TupleExprContext::getRuleIndex() const {
  return bapelParser::RuleTupleExpr;
}


antlrcpp::Any bapelParser::TupleExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTupleExpr(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TupleExprContext* bapelParser::tupleExpr() {
  TupleExprContext *_localctx = _tracker.createInstance<TupleExprContext>(_ctx, getState());
  enterRule(_localctx, 162, bapelParser::RuleTupleExpr);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(868);
    _errHandler->sync(this);
    switch (getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 92, _ctx)) {
    case 1: {
      enterOuterAlt(_localctx, 1);
      setState(862);
      match(bapelParser::LPAREN);
      setState(863);
      match(bapelParser::RPAREN);
      break;
    }

    case 2: {
      enterOuterAlt(_localctx, 2);
      setState(864);
      match(bapelParser::LPAREN);
      setState(865);
      tupleExprArgs();
      setState(866);
      match(bapelParser::RPAREN);
      break;
    }

    default:
      break;
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- TupleExprArgsContext ------------------------------------------------------------------

bapelParser::TupleExprArgsContext::TupleExprArgsContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

std::vector<bapelParser::ExpressionContext *> bapelParser::TupleExprArgsContext::expression() {
  return getRuleContexts<bapelParser::ExpressionContext>();
}

bapelParser::ExpressionContext* bapelParser::TupleExprArgsContext::expression(size_t i) {
  return getRuleContext<bapelParser::ExpressionContext>(i);
}

std::vector<tree::TerminalNode *> bapelParser::TupleExprArgsContext::COMMA() {
  return getTokens(bapelParser::COMMA);
}

tree::TerminalNode* bapelParser::TupleExprArgsContext::COMMA(size_t i) {
  return getToken(bapelParser::COMMA, i);
}


size_t bapelParser::TupleExprArgsContext::getRuleIndex() const {
  return bapelParser::RuleTupleExprArgs;
}


antlrcpp::Any bapelParser::TupleExprArgsContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitTupleExprArgs(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::TupleExprArgsContext* bapelParser::tupleExprArgs() {
  TupleExprArgsContext *_localctx = _tracker.createInstance<TupleExprArgsContext>(_ctx, getState());
  enterRule(_localctx, 164, bapelParser::RuleTupleExprArgs);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(870);
    expression();
    setState(873); 
    _errHandler->sync(this);
    _la = _input->LA(1);
    do {
      setState(871);
      match(bapelParser::COMMA);
      setState(872);
      expression();
      setState(875); 
      _errHandler->sync(this);
      _la = _input->LA(1);
    } while (_la == bapelParser::COMMA);
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- VarExprContext ------------------------------------------------------------------

bapelParser::VarExprContext::VarExprContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::IdContext* bapelParser::VarExprContext::id() {
  return getRuleContext<bapelParser::IdContext>(0);
}


size_t bapelParser::VarExprContext::getRuleIndex() const {
  return bapelParser::RuleVarExpr;
}


antlrcpp::Any bapelParser::VarExprContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitVarExpr(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::VarExprContext* bapelParser::varExpr() {
  VarExprContext *_localctx = _tracker.createInstance<VarExprContext>(_ctx, getState());
  enterRule(_localctx, 166, bapelParser::RuleVarExpr);

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    enterOuterAlt(_localctx, 1);
    setState(877);
    id();
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- IdContext ------------------------------------------------------------------

bapelParser::IdContext::IdContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

bapelParser::IdTokensContext* bapelParser::IdContext::idTokens() {
  return getRuleContext<bapelParser::IdTokensContext>(0);
}

tree::TerminalNode* bapelParser::IdContext::LPAREN() {
  return getToken(bapelParser::LPAREN, 0);
}

tree::TerminalNode* bapelParser::IdContext::RPAREN() {
  return getToken(bapelParser::RPAREN, 0);
}

tree::TerminalNode* bapelParser::IdContext::OR() {
  return getToken(bapelParser::OR, 0);
}

tree::TerminalNode* bapelParser::IdContext::AND() {
  return getToken(bapelParser::AND, 0);
}

tree::TerminalNode* bapelParser::IdContext::NE() {
  return getToken(bapelParser::NE, 0);
}

tree::TerminalNode* bapelParser::IdContext::EQ() {
  return getToken(bapelParser::EQ, 0);
}

tree::TerminalNode* bapelParser::IdContext::GT() {
  return getToken(bapelParser::GT, 0);
}

tree::TerminalNode* bapelParser::IdContext::GE() {
  return getToken(bapelParser::GE, 0);
}

tree::TerminalNode* bapelParser::IdContext::LT() {
  return getToken(bapelParser::LT, 0);
}

tree::TerminalNode* bapelParser::IdContext::LE() {
  return getToken(bapelParser::LE, 0);
}

tree::TerminalNode* bapelParser::IdContext::PLUS() {
  return getToken(bapelParser::PLUS, 0);
}

tree::TerminalNode* bapelParser::IdContext::MINUS() {
  return getToken(bapelParser::MINUS, 0);
}

tree::TerminalNode* bapelParser::IdContext::MUL() {
  return getToken(bapelParser::MUL, 0);
}

tree::TerminalNode* bapelParser::IdContext::DIV() {
  return getToken(bapelParser::DIV, 0);
}

tree::TerminalNode* bapelParser::IdContext::NOT() {
  return getToken(bapelParser::NOT, 0);
}


size_t bapelParser::IdContext::getRuleIndex() const {
  return bapelParser::RuleId;
}


antlrcpp::Any bapelParser::IdContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitId(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::IdContext* bapelParser::id() {
  IdContext *_localctx = _tracker.createInstance<IdContext>(_ctx, getState());
  enterRule(_localctx, 168, bapelParser::RuleId);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    setState(883);
    _errHandler->sync(this);
    switch (_input->LA(1)) {
      case bapelParser::IDENTIFIER: {
        enterOuterAlt(_localctx, 1);
        setState(879);
        idTokens();
        break;
      }

      case bapelParser::LPAREN: {
        enterOuterAlt(_localctx, 2);
        setState(880);
        match(bapelParser::LPAREN);
        setState(881);
        _la = _input->LA(1);
        if (!((((_la & ~ 0x3fULL) == 0) &&
          ((1ULL << _la) & ((1ULL << bapelParser::OR)
          | (1ULL << bapelParser::AND)
          | (1ULL << bapelParser::NE)
          | (1ULL << bapelParser::EQ)
          | (1ULL << bapelParser::GE)
          | (1ULL << bapelParser::LE)
          | (1ULL << bapelParser::GT)
          | (1ULL << bapelParser::LT)
          | (1ULL << bapelParser::PLUS)
          | (1ULL << bapelParser::MINUS)
          | (1ULL << bapelParser::MUL)
          | (1ULL << bapelParser::DIV)
          | (1ULL << bapelParser::NOT))) != 0))) {
        _errHandler->recoverInline(this);
        }
        else {
          _errHandler->reportMatch(this);
          consume();
        }
        setState(882);
        match(bapelParser::RPAREN);
        break;
      }

    default:
      throw NoViableAltException(this);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

//----------------- IdTokensContext ------------------------------------------------------------------

bapelParser::IdTokensContext::IdTokensContext(ParserRuleContext *parent, size_t invokingState)
  : ParserRuleContext(parent, invokingState) {
}

std::vector<tree::TerminalNode *> bapelParser::IdTokensContext::IDENTIFIER() {
  return getTokens(bapelParser::IDENTIFIER);
}

tree::TerminalNode* bapelParser::IdTokensContext::IDENTIFIER(size_t i) {
  return getToken(bapelParser::IDENTIFIER, i);
}

std::vector<tree::TerminalNode *> bapelParser::IdTokensContext::DOUBLE_COLON() {
  return getTokens(bapelParser::DOUBLE_COLON);
}

tree::TerminalNode* bapelParser::IdTokensContext::DOUBLE_COLON(size_t i) {
  return getToken(bapelParser::DOUBLE_COLON, i);
}

std::vector<tree::TerminalNode *> bapelParser::IdTokensContext::SET() {
  return getTokens(bapelParser::SET);
}

tree::TerminalNode* bapelParser::IdTokensContext::SET(size_t i) {
  return getToken(bapelParser::SET, i);
}


size_t bapelParser::IdTokensContext::getRuleIndex() const {
  return bapelParser::RuleIdTokens;
}


antlrcpp::Any bapelParser::IdTokensContext::accept(tree::ParseTreeVisitor *visitor) {
  if (auto parserVisitor = dynamic_cast<bapelVisitor*>(visitor))
    return parserVisitor->visitIdTokens(this);
  else
    return visitor->visitChildren(this);
}

bapelParser::IdTokensContext* bapelParser::idTokens() {
  IdTokensContext *_localctx = _tracker.createInstance<IdTokensContext>(_ctx, getState());
  enterRule(_localctx, 170, bapelParser::RuleIdTokens);
  size_t _la = 0;

#if __cplusplus > 201703L
  auto onExit = finally([=, this] {
#else
  auto onExit = finally([=] {
#endif
    exitRule();
  });
  try {
    size_t alt;
    enterOuterAlt(_localctx, 1);
    setState(885);
    match(bapelParser::IDENTIFIER);
    setState(890);
    _errHandler->sync(this);
    alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 95, _ctx);
    while (alt != 2 && alt != atn::ATN::INVALID_ALT_NUMBER) {
      if (alt == 1) {
        setState(886);
        match(bapelParser::DOUBLE_COLON);
        setState(887);
        _la = _input->LA(1);
        if (!(_la == bapelParser::SET

        || _la == bapelParser::IDENTIFIER)) {
        _errHandler->recoverInline(this);
        }
        else {
          _errHandler->reportMatch(this);
          consume();
        } 
      }
      setState(892);
      _errHandler->sync(this);
      alt = getInterpreter<atn::ParserATNSimulator>()->adaptivePredict(_input, 95, _ctx);
    }
   
  }
  catch (RecognitionException &e) {
    _errHandler->reportError(this, e);
    _localctx->exception = std::current_exception();
    _errHandler->recover(this, _localctx->exception);
  }

  return _localctx;
}

bool bapelParser::sempred(RuleContext *context, size_t ruleIndex, size_t predicateIndex) {
  switch (ruleIndex) {
    case 33: return appTypeSempred(dynamic_cast<AppTypeContext *>(context), predicateIndex);
    case 63: return logicalOrExprSempred(dynamic_cast<LogicalOrExprContext *>(context), predicateIndex);
    case 64: return logicalAndExprSempred(dynamic_cast<LogicalAndExprContext *>(context), predicateIndex);
    case 65: return equalityExprSempred(dynamic_cast<EqualityExprContext *>(context), predicateIndex);
    case 66: return comparisonExprSempred(dynamic_cast<ComparisonExprContext *>(context), predicateIndex);
    case 67: return additiveExprSempred(dynamic_cast<AdditiveExprContext *>(context), predicateIndex);
    case 68: return multiplicativeExprSempred(dynamic_cast<MultiplicativeExprContext *>(context), predicateIndex);
    case 70: return applicativeExprSempred(dynamic_cast<ApplicativeExprContext *>(context), predicateIndex);
    case 75: return projectionExprSempred(dynamic_cast<ProjectionExprContext *>(context), predicateIndex);

  default:
    break;
  }
  return true;
}

bool bapelParser::appTypeSempred(AppTypeContext *_localctx, size_t predicateIndex) {
  switch (predicateIndex) {
    case 0: return precpred(_ctx, 2);

  default:
    break;
  }
  return true;
}

bool bapelParser::logicalOrExprSempred(LogicalOrExprContext *_localctx, size_t predicateIndex) {
  switch (predicateIndex) {
    case 1: return precpred(_ctx, 2);

  default:
    break;
  }
  return true;
}

bool bapelParser::logicalAndExprSempred(LogicalAndExprContext *_localctx, size_t predicateIndex) {
  switch (predicateIndex) {
    case 2: return precpred(_ctx, 2);

  default:
    break;
  }
  return true;
}

bool bapelParser::equalityExprSempred(EqualityExprContext *_localctx, size_t predicateIndex) {
  switch (predicateIndex) {
    case 3: return precpred(_ctx, 2);

  default:
    break;
  }
  return true;
}

bool bapelParser::comparisonExprSempred(ComparisonExprContext *_localctx, size_t predicateIndex) {
  switch (predicateIndex) {
    case 4: return precpred(_ctx, 2);

  default:
    break;
  }
  return true;
}

bool bapelParser::additiveExprSempred(AdditiveExprContext *_localctx, size_t predicateIndex) {
  switch (predicateIndex) {
    case 5: return precpred(_ctx, 2);

  default:
    break;
  }
  return true;
}

bool bapelParser::multiplicativeExprSempred(MultiplicativeExprContext *_localctx, size_t predicateIndex) {
  switch (predicateIndex) {
    case 6: return precpred(_ctx, 2);

  default:
    break;
  }
  return true;
}

bool bapelParser::applicativeExprSempred(ApplicativeExprContext *_localctx, size_t predicateIndex) {
  switch (predicateIndex) {
    case 7: return precpred(_ctx, 2);

  default:
    break;
  }
  return true;
}

bool bapelParser::projectionExprSempred(ProjectionExprContext *_localctx, size_t predicateIndex) {
  switch (predicateIndex) {
    case 8: return precpred(_ctx, 3);
    case 9: return precpred(_ctx, 2);

  default:
    break;
  }
  return true;
}

// Static vars and initialization.
std::vector<dfa::DFA> bapelParser::_decisionToDFA;
atn::PredictionContextCache bapelParser::_sharedContextCache;

// We own the ATN which in turn owns the ATN states.
atn::ATN bapelParser::_atn;
std::vector<uint16_t> bapelParser::_serializedATN;

std::vector<std::string> bapelParser::_ruleNames = {
  "sourceFile", "moduleHeader", "implementsHeader", "workspace", "packagesSection", 
  "packageRule", "importsSection", "implsSection", "flagsSection", "moduleID", 
  "filename", "sources", "source", "traitDecl", "traitMethod", "implBlock", 
  "declNoExport", "functionNoExport", "functionArgs", "arg", "decl", "unexportedDecl", 
  "declNoTerm", "termDecl", "typeDecl", "typeAbstraction", "boundedTvar", 
  "tvar", "traitBound", "type_", "forallType", "functionType", "ptrType", 
  "appType", "primaryType", "arrayType", "structType", "fields", "field", 
  "tupleType", "tupleTypeArgs", "variantType", "tags", "tag", "expression", 
  "expressionWithoutBlock", "expressionWithBlock", "assignTerm", "returnTerm", 
  "ifTerm", "forTerm", "lambdaTerm", "matchTerm", "matchArms", "matchArm", 
  "setTerm", "blockExpr", "blockStatements", "statements", "statement", 
  "letStatement", "expressionStatement", "operatorExpr", "logicalOrExpr", 
  "logicalAndExpr", "equalityExpr", "comparisonExpr", "additiveExpr", "multiplicativeExpr", 
  "unaryExpr", "applicativeExpr", "typeApplicativeExpr", "typeApplicativeArgs", 
  "basePrimaryExpr", "primaryExpr", "projectionExpr", "derefExpr", "injectionExpr", 
  "structExpr", "labelValues", "labelValue", "tupleExpr", "tupleExprArgs", 
  "varExpr", "id", "idTokens"
};

std::vector<std::string> bapelParser::_literalNames = {
  "", "'workspace'", "'packages'", "'prefix'", "'module'", "'implements'", 
  "'imports'", "'impls'", "'flags'", "'in'", "'pub'", "'decl'", "'fn'", 
  "'type'", "'forall'", "'struct'", "'variant'", "'match'", "'set'", "'let'", 
  "'return'", "'if'", "'else'", "'for'", "'trait'", "'impl'", "'::'", "'->'", 
  "'=>'", "'<-'", "'||'", "'&&'", "'!='", "'=='", "'>='", "'<='", "'>'", 
  "'<'", "'+'", "'-'", "'*'", "'/'", "'!'", "'&'", "'.'", "'='", "','", 
  "';'", "':'", "'{'", "'}'", "'('", "')'", "'['", "']'", "'''"
};

std::vector<std::string> bapelParser::_symbolicNames = {
  "", "WORKSPACE", "PACKAGES", "PREFIX", "MODULE", "IMPLEMENTS", "IMPORTS", 
  "IMPLS", "FLAGS", "IN", "PUB", "DECL", "FN", "TYPE", "FORALL", "STRUCT", 
  "VARIANT", "MATCH", "SET", "LET", "RETURN", "IF", "ELSE", "FOR", "TRAIT", 
  "IMPL", "DOUBLE_COLON", "ARROW", "FAT_ARROW", "LARROW", "OR", "AND", "NE", 
  "EQ", "GE", "LE", "GT", "LT", "PLUS", "MINUS", "MUL", "DIV", "NOT", "AMP", 
  "DOT", "ASSIGN", "COMMA", "SEMI", "COLON", "LBRACE", "RBRACE", "LPAREN", 
  "RPAREN", "LBRACKET", "RBRACKET", "SINGLE_QUOTE", "IDENTIFIER", "INT_LITERAL", 
  "FLOAT_LITERAL", "RUNE_LITERAL", "STRING_LITERAL", "RAW_STRING_LITERAL", 
  "UNTERMINATED_RUNE_LITERAL", "UNTERMINATED_STRING_LITERAL", "UNTERMINATED_RAW_STRING_LITERAL", 
  "UNTERMINATED_BLOCK_COMMENT", "WS", "LINE_COMMENT", "BLOCK_COMMENT"
};

dfa::Vocabulary bapelParser::_vocabulary(_literalNames, _symbolicNames);

std::vector<std::string> bapelParser::_tokenNames;

bapelParser::Initializer::Initializer() {
	for (size_t i = 0; i < _symbolicNames.size(); ++i) {
		std::string name = _vocabulary.getLiteralName(i);
		if (name.empty()) {
			name = _vocabulary.getSymbolicName(i);
		}

		if (name.empty()) {
			_tokenNames.push_back("<INVALID>");
		} else {
      _tokenNames.push_back(name);
    }
	}

  static const uint16_t serializedATNSegment0[] = {
    0x3, 0x608b, 0xa72a, 0x8133, 0xb9ed, 0x417c, 0x3be7, 0x7786, 0x5964, 
       0x3, 0x46, 0x380, 0x4, 0x2, 0x9, 0x2, 0x4, 0x3, 0x9, 0x3, 0x4, 0x4, 
       0x9, 0x4, 0x4, 0x5, 0x9, 0x5, 0x4, 0x6, 0x9, 0x6, 0x4, 0x7, 0x9, 
       0x7, 0x4, 0x8, 0x9, 0x8, 0x4, 0x9, 0x9, 0x9, 0x4, 0xa, 0x9, 0xa, 
       0x4, 0xb, 0x9, 0xb, 0x4, 0xc, 0x9, 0xc, 0x4, 0xd, 0x9, 0xd, 0x4, 
       0xe, 0x9, 0xe, 0x4, 0xf, 0x9, 0xf, 0x4, 0x10, 0x9, 0x10, 0x4, 0x11, 
       0x9, 0x11, 0x4, 0x12, 0x9, 0x12, 0x4, 0x13, 0x9, 0x13, 0x4, 0x14, 
       0x9, 0x14, 0x4, 0x15, 0x9, 0x15, 0x4, 0x16, 0x9, 0x16, 0x4, 0x17, 
       0x9, 0x17, 0x4, 0x18, 0x9, 0x18, 0x4, 0x19, 0x9, 0x19, 0x4, 0x1a, 
       0x9, 0x1a, 0x4, 0x1b, 0x9, 0x1b, 0x4, 0x1c, 0x9, 0x1c, 0x4, 0x1d, 
       0x9, 0x1d, 0x4, 0x1e, 0x9, 0x1e, 0x4, 0x1f, 0x9, 0x1f, 0x4, 0x20, 
       0x9, 0x20, 0x4, 0x21, 0x9, 0x21, 0x4, 0x22, 0x9, 0x22, 0x4, 0x23, 
       0x9, 0x23, 0x4, 0x24, 0x9, 0x24, 0x4, 0x25, 0x9, 0x25, 0x4, 0x26, 
       0x9, 0x26, 0x4, 0x27, 0x9, 0x27, 0x4, 0x28, 0x9, 0x28, 0x4, 0x29, 
       0x9, 0x29, 0x4, 0x2a, 0x9, 0x2a, 0x4, 0x2b, 0x9, 0x2b, 0x4, 0x2c, 
       0x9, 0x2c, 0x4, 0x2d, 0x9, 0x2d, 0x4, 0x2e, 0x9, 0x2e, 0x4, 0x2f, 
       0x9, 0x2f, 0x4, 0x30, 0x9, 0x30, 0x4, 0x31, 0x9, 0x31, 0x4, 0x32, 
       0x9, 0x32, 0x4, 0x33, 0x9, 0x33, 0x4, 0x34, 0x9, 0x34, 0x4, 0x35, 
       0x9, 0x35, 0x4, 0x36, 0x9, 0x36, 0x4, 0x37, 0x9, 0x37, 0x4, 0x38, 
       0x9, 0x38, 0x4, 0x39, 0x9, 0x39, 0x4, 0x3a, 0x9, 0x3a, 0x4, 0x3b, 
       0x9, 0x3b, 0x4, 0x3c, 0x9, 0x3c, 0x4, 0x3d, 0x9, 0x3d, 0x4, 0x3e, 
       0x9, 0x3e, 0x4, 0x3f, 0x9, 0x3f, 0x4, 0x40, 0x9, 0x40, 0x4, 0x41, 
       0x9, 0x41, 0x4, 0x42, 0x9, 0x42, 0x4, 0x43, 0x9, 0x43, 0x4, 0x44, 
       0x9, 0x44, 0x4, 0x45, 0x9, 0x45, 0x4, 0x46, 0x9, 0x46, 0x4, 0x47, 
       0x9, 0x47, 0x4, 0x48, 0x9, 0x48, 0x4, 0x49, 0x9, 0x49, 0x4, 0x4a, 
       0x9, 0x4a, 0x4, 0x4b, 0x9, 0x4b, 0x4, 0x4c, 0x9, 0x4c, 0x4, 0x4d, 
       0x9, 0x4d, 0x4, 0x4e, 0x9, 0x4e, 0x4, 0x4f, 0x9, 0x4f, 0x4, 0x50, 
       0x9, 0x50, 0x4, 0x51, 0x9, 0x51, 0x4, 0x52, 0x9, 0x52, 0x4, 0x53, 
       0x9, 0x53, 0x4, 0x54, 0x9, 0x54, 0x4, 0x55, 0x9, 0x55, 0x4, 0x56, 
       0x9, 0x56, 0x4, 0x57, 0x9, 0x57, 0x3, 0x2, 0x3, 0x2, 0x5, 0x2, 0xb1, 
       0xa, 0x2, 0x3, 0x2, 0x5, 0x2, 0xb4, 0xa, 0x2, 0x3, 0x2, 0x5, 0x2, 
       0xb7, 0xa, 0x2, 0x3, 0x2, 0x5, 0x2, 0xba, 0xa, 0x2, 0x3, 0x2, 0x3, 
       0x2, 0x3, 0x2, 0x3, 0x2, 0x5, 0x2, 0xc0, 0xa, 0x2, 0x3, 0x2, 0x5, 
       0x2, 0xc3, 0xa, 0x2, 0x3, 0x2, 0x5, 0x2, 0xc6, 0xa, 0x2, 0x3, 0x2, 
       0x5, 0x2, 0xc9, 0xa, 0x2, 0x3, 0x2, 0x3, 0x2, 0x5, 0x2, 0xcd, 0xa, 
       0x2, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x4, 0x3, 0x4, 0x3, 0x4, 
       0x3, 0x5, 0x3, 0x5, 0x3, 0x5, 0x3, 0x5, 0x3, 0x5, 0x3, 0x6, 0x3, 
       0x6, 0x3, 0x6, 0x6, 0x6, 0xdd, 0xa, 0x6, 0xd, 0x6, 0xe, 0x6, 0xde, 
       0x3, 0x6, 0x3, 0x6, 0x3, 0x7, 0x3, 0x7, 0x3, 0x7, 0x3, 0x7, 0x3, 
       0x7, 0x3, 0x7, 0x3, 0x7, 0x3, 0x7, 0x3, 0x7, 0x3, 0x7, 0x5, 0x7, 
       0xed, 0xa, 0x7, 0x3, 0x8, 0x3, 0x8, 0x3, 0x8, 0x6, 0x8, 0xf2, 0xa, 
       0x8, 0xd, 0x8, 0xe, 0x8, 0xf3, 0x3, 0x8, 0x3, 0x8, 0x3, 0x9, 0x3, 
       0x9, 0x3, 0x9, 0x6, 0x9, 0xfb, 0xa, 0x9, 0xd, 0x9, 0xe, 0x9, 0xfc, 
       0x3, 0x9, 0x3, 0x9, 0x3, 0xa, 0x3, 0xa, 0x3, 0xa, 0x6, 0xa, 0x104, 
       0xa, 0xa, 0xd, 0xa, 0xe, 0xa, 0x105, 0x3, 0xa, 0x3, 0xa, 0x3, 0xb, 
       0x3, 0xb, 0x3, 0xb, 0x7, 0xb, 0x10d, 0xa, 0xb, 0xc, 0xb, 0xe, 0xb, 
       0x110, 0xb, 0xb, 0x3, 0xc, 0x3, 0xc, 0x3, 0xd, 0x6, 0xd, 0x115, 0xa, 
       0xd, 0xd, 0xd, 0xe, 0xd, 0x116, 0x3, 0xe, 0x3, 0xe, 0x3, 0xe, 0x3, 
       0xe, 0x3, 0xe, 0x3, 0xe, 0x3, 0xe, 0x3, 0xe, 0x5, 0xe, 0x121, 0xa, 
       0xe, 0x3, 0xf, 0x3, 0xf, 0x3, 0xf, 0x5, 0xf, 0x126, 0xa, 0xf, 0x3, 
       0xf, 0x3, 0xf, 0x7, 0xf, 0x12a, 0xa, 0xf, 0xc, 0xf, 0xe, 0xf, 0x12d, 
       0xb, 0xf, 0x3, 0xf, 0x3, 0xf, 0x3, 0x10, 0x3, 0x10, 0x3, 0x10, 0x3, 
       0x10, 0x3, 0x10, 0x3, 0x10, 0x3, 0x11, 0x3, 0x11, 0x5, 0x11, 0x139, 
       0xa, 0x11, 0x3, 0x11, 0x3, 0x11, 0x3, 0x11, 0x3, 0x11, 0x3, 0x11, 
       0x7, 0x11, 0x140, 0xa, 0x11, 0xc, 0x11, 0xe, 0x11, 0x143, 0xb, 0x11, 
       0x3, 0x11, 0x3, 0x11, 0x3, 0x11, 0x3, 0x11, 0x5, 0x11, 0x149, 0xa, 
       0x11, 0x3, 0x11, 0x3, 0x11, 0x3, 0x11, 0x7, 0x11, 0x14e, 0xa, 0x11, 
       0xc, 0x11, 0xe, 0x11, 0x151, 0xb, 0x11, 0x3, 0x11, 0x3, 0x11, 0x5, 
       0x11, 0x155, 0xa, 0x11, 0x3, 0x12, 0x3, 0x12, 0x3, 0x12, 0x3, 0x12, 
       0x3, 0x12, 0x5, 0x12, 0x15c, 0xa, 0x12, 0x3, 0x13, 0x3, 0x13, 0x3, 
       0x13, 0x5, 0x13, 0x161, 0xa, 0x13, 0x3, 0x13, 0x3, 0x13, 0x3, 0x13, 
       0x3, 0x13, 0x3, 0x13, 0x3, 0x14, 0x3, 0x14, 0x3, 0x14, 0x3, 0x14, 
       0x7, 0x14, 0x16c, 0xa, 0x14, 0xc, 0x14, 0xe, 0x14, 0x16f, 0xb, 0x14, 
       0x5, 0x14, 0x171, 0xa, 0x14, 0x3, 0x14, 0x3, 0x14, 0x3, 0x15, 0x3, 
       0x15, 0x3, 0x15, 0x3, 0x15, 0x3, 0x16, 0x3, 0x16, 0x3, 0x16, 0x5, 
       0x16, 0x17c, 0xa, 0x16, 0x3, 0x17, 0x3, 0x17, 0x5, 0x17, 0x180, 0xa, 
       0x17, 0x3, 0x18, 0x3, 0x18, 0x3, 0x18, 0x5, 0x18, 0x185, 0xa, 0x18, 
       0x3, 0x19, 0x3, 0x19, 0x3, 0x19, 0x3, 0x19, 0x3, 0x1a, 0x3, 0x1a, 
       0x3, 0x1a, 0x5, 0x1a, 0x18e, 0xa, 0x1a, 0x3, 0x1a, 0x3, 0x1a, 0x3, 
       0x1a, 0x3, 0x1a, 0x3, 0x1a, 0x3, 0x1a, 0x5, 0x1a, 0x196, 0xa, 0x1a, 
       0x5, 0x1a, 0x198, 0xa, 0x1a, 0x3, 0x1b, 0x3, 0x1b, 0x3, 0x1b, 0x3, 
       0x1b, 0x7, 0x1b, 0x19e, 0xa, 0x1b, 0xc, 0x1b, 0xe, 0x1b, 0x1a1, 0xb, 
       0x1b, 0x3, 0x1b, 0x3, 0x1b, 0x3, 0x1c, 0x3, 0x1c, 0x3, 0x1c, 0x5, 
       0x1c, 0x1a8, 0xa, 0x1c, 0x3, 0x1d, 0x3, 0x1d, 0x3, 0x1d, 0x3, 0x1e, 
       0x3, 0x1e, 0x3, 0x1e, 0x7, 0x1e, 0x1b0, 0xa, 0x1e, 0xc, 0x1e, 0xe, 
       0x1e, 0x1b3, 0xb, 0x1e, 0x3, 0x1f, 0x3, 0x1f, 0x3, 0x20, 0x3, 0x20, 
       0x3, 0x20, 0x3, 0x20, 0x3, 0x20, 0x5, 0x20, 0x1bc, 0xa, 0x20, 0x3, 
       0x21, 0x3, 0x21, 0x3, 0x21, 0x5, 0x21, 0x1c1, 0xa, 0x21, 0x3, 0x22, 
       0x3, 0x22, 0x3, 0x22, 0x5, 0x22, 0x1c6, 0xa, 0x22, 0x3, 0x23, 0x3, 
       0x23, 0x3, 0x23, 0x3, 0x23, 0x3, 0x23, 0x7, 0x23, 0x1cd, 0xa, 0x23, 
       0xc, 0x23, 0xe, 0x23, 0x1d0, 0xb, 0x23, 0x3, 0x24, 0x3, 0x24, 0x3, 
       0x24, 0x3, 0x24, 0x3, 0x24, 0x3, 0x24, 0x3, 0x24, 0x3, 0x24, 0x3, 
       0x24, 0x3, 0x24, 0x3, 0x24, 0x5, 0x24, 0x1dd, 0xa, 0x24, 0x3, 0x25, 
       0x3, 0x25, 0x3, 0x25, 0x3, 0x25, 0x3, 0x25, 0x3, 0x25, 0x3, 0x25, 
       0x3, 0x25, 0x3, 0x25, 0x3, 0x25, 0x3, 0x25, 0x3, 0x25, 0x3, 0x25, 
       0x3, 0x25, 0x3, 0x25, 0x3, 0x25, 0x3, 0x25, 0x5, 0x25, 0x1f0, 0xa, 
       0x25, 0x3, 0x26, 0x3, 0x26, 0x3, 0x26, 0x3, 0x26, 0x5, 0x26, 0x1f6, 
       0xa, 0x26, 0x5, 0x26, 0x1f8, 0xa, 0x26, 0x3, 0x26, 0x3, 0x26, 0x3, 
       0x27, 0x3, 0x27, 0x3, 0x27, 0x7, 0x27, 0x1ff, 0xa, 0x27, 0xc, 0x27, 
       0xe, 0x27, 0x202, 0xb, 0x27, 0x3, 0x28, 0x3, 0x28, 0x3, 0x28, 0x3, 
       0x28, 0x3, 0x29, 0x3, 0x29, 0x3, 0x29, 0x3, 0x29, 0x3, 0x29, 0x3, 
       0x29, 0x5, 0x29, 0x20e, 0xa, 0x29, 0x3, 0x2a, 0x3, 0x2a, 0x3, 0x2a, 
       0x6, 0x2a, 0x213, 0xa, 0x2a, 0xd, 0x2a, 0xe, 0x2a, 0x214, 0x3, 0x2b, 
       0x3, 0x2b, 0x3, 0x2b, 0x3, 0x2b, 0x5, 0x2b, 0x21b, 0xa, 0x2b, 0x5, 
       0x2b, 0x21d, 0xa, 0x2b, 0x3, 0x2b, 0x3, 0x2b, 0x3, 0x2c, 0x3, 0x2c, 
       0x3, 0x2c, 0x7, 0x2c, 0x224, 0xa, 0x2c, 0xc, 0x2c, 0xe, 0x2c, 0x227, 
       0xb, 0x2c, 0x3, 0x2d, 0x3, 0x2d, 0x3, 0x2d, 0x3, 0x2e, 0x3, 0x2e, 
       0x5, 0x2e, 0x22e, 0xa, 0x2e, 0x3, 0x2f, 0x3, 0x2f, 0x3, 0x2f, 0x5, 
       0x2f, 0x233, 0xa, 0x2f, 0x3, 0x30, 0x3, 0x30, 0x3, 0x30, 0x3, 0x30, 
       0x3, 0x30, 0x3, 0x30, 0x5, 0x30, 0x23b, 0xa, 0x30, 0x3, 0x31, 0x3, 
       0x31, 0x5, 0x31, 0x23f, 0xa, 0x31, 0x3, 0x31, 0x3, 0x31, 0x3, 0x31, 
       0x3, 0x32, 0x3, 0x32, 0x3, 0x32, 0x3, 0x33, 0x3, 0x33, 0x3, 0x33, 
       0x3, 0x33, 0x3, 0x33, 0x3, 0x33, 0x5, 0x33, 0x24d, 0xa, 0x33, 0x5, 
       0x33, 0x24f, 0xa, 0x33, 0x3, 0x34, 0x3, 0x34, 0x3, 0x34, 0x3, 0x34, 
       0x3, 0x35, 0x3, 0x35, 0x5, 0x35, 0x257, 0xa, 0x35, 0x3, 0x35, 0x3, 
       0x35, 0x3, 0x35, 0x3, 0x36, 0x3, 0x36, 0x3, 0x36, 0x3, 0x36, 0x3, 
       0x36, 0x5, 0x36, 0x261, 0xa, 0x36, 0x3, 0x36, 0x3, 0x36, 0x3, 0x37, 
       0x3, 0x37, 0x3, 0x37, 0x7, 0x37, 0x268, 0xa, 0x37, 0xc, 0x37, 0xe, 
       0x37, 0x26b, 0xb, 0x37, 0x3, 0x38, 0x3, 0x38, 0x3, 0x38, 0x3, 0x38, 
       0x3, 0x38, 0x3, 0x39, 0x3, 0x39, 0x3, 0x39, 0x3, 0x39, 0x3, 0x39, 
       0x5, 0x39, 0x277, 0xa, 0x39, 0x3, 0x39, 0x3, 0x39, 0x3, 0x3a, 0x3, 
       0x3a, 0x3, 0x3a, 0x3, 0x3a, 0x3, 0x3b, 0x3, 0x3b, 0x5, 0x3b, 0x281, 
       0xa, 0x3b, 0x3, 0x3b, 0x5, 0x3b, 0x284, 0xa, 0x3b, 0x3, 0x3c, 0x6, 
       0x3c, 0x287, 0xa, 0x3c, 0xd, 0x3c, 0xe, 0x3c, 0x288, 0x3, 0x3d, 0x3, 
       0x3d, 0x5, 0x3d, 0x28d, 0xa, 0x3d, 0x3, 0x3e, 0x3, 0x3e, 0x3, 0x3e, 
       0x3, 0x3e, 0x3, 0x3e, 0x3, 0x3e, 0x3, 0x3e, 0x3, 0x3e, 0x3, 0x3e, 
       0x3, 0x3e, 0x3, 0x3e, 0x3, 0x3e, 0x3, 0x3e, 0x3, 0x3e, 0x5, 0x3e, 
       0x29d, 0xa, 0x3e, 0x3, 0x3f, 0x3, 0x3f, 0x3, 0x3f, 0x3, 0x3f, 0x3, 
       0x3f, 0x5, 0x3f, 0x2a4, 0xa, 0x3f, 0x5, 0x3f, 0x2a6, 0xa, 0x3f, 0x3, 
       0x40, 0x3, 0x40, 0x3, 0x41, 0x3, 0x41, 0x3, 0x41, 0x3, 0x41, 0x3, 
       0x41, 0x3, 0x41, 0x7, 0x41, 0x2b0, 0xa, 0x41, 0xc, 0x41, 0xe, 0x41, 
       0x2b3, 0xb, 0x41, 0x3, 0x42, 0x3, 0x42, 0x3, 0x42, 0x3, 0x42, 0x3, 
       0x42, 0x3, 0x42, 0x7, 0x42, 0x2bb, 0xa, 0x42, 0xc, 0x42, 0xe, 0x42, 
       0x2be, 0xb, 0x42, 0x3, 0x43, 0x3, 0x43, 0x3, 0x43, 0x3, 0x43, 0x3, 
       0x43, 0x3, 0x43, 0x5, 0x43, 0x2c6, 0xa, 0x43, 0x3, 0x43, 0x7, 0x43, 
       0x2c9, 0xa, 0x43, 0xc, 0x43, 0xe, 0x43, 0x2cc, 0xb, 0x43, 0x3, 0x44, 
       0x3, 0x44, 0x3, 0x44, 0x3, 0x44, 0x3, 0x44, 0x3, 0x44, 0x5, 0x44, 
       0x2d4, 0xa, 0x44, 0x3, 0x44, 0x7, 0x44, 0x2d7, 0xa, 0x44, 0xc, 0x44, 
       0xe, 0x44, 0x2da, 0xb, 0x44, 0x3, 0x45, 0x3, 0x45, 0x3, 0x45, 0x3, 
       0x45, 0x3, 0x45, 0x3, 0x45, 0x5, 0x45, 0x2e2, 0xa, 0x45, 0x3, 0x45, 
       0x7, 0x45, 0x2e5, 0xa, 0x45, 0xc, 0x45, 0xe, 0x45, 0x2e8, 0xb, 0x45, 
       0x3, 0x46, 0x3, 0x46, 0x3, 0x46, 0x3, 0x46, 0x3, 0x46, 0x3, 0x46, 
       0x5, 0x46, 0x2f0, 0xa, 0x46, 0x3, 0x46, 0x7, 0x46, 0x2f3, 0xa, 0x46, 
       0xc, 0x46, 0xe, 0x46, 0x2f6, 0xb, 0x46, 0x3, 0x47, 0x3, 0x47, 0x5, 
       0x47, 0x2fa, 0xa, 0x47, 0x3, 0x47, 0x3, 0x47, 0x5, 0x47, 0x2fe, 0xa, 
       0x47, 0x3, 0x48, 0x3, 0x48, 0x3, 0x48, 0x3, 0x48, 0x3, 0x48, 0x7, 
       0x48, 0x305, 0xa, 0x48, 0xc, 0x48, 0xe, 0x48, 0x308, 0xb, 0x48, 0x3, 
       0x49, 0x3, 0x49, 0x3, 0x49, 0x3, 0x49, 0x5, 0x49, 0x30e, 0xa, 0x49, 
       0x3, 0x4a, 0x3, 0x4a, 0x3, 0x4a, 0x3, 0x4a, 0x3, 0x4a, 0x3, 0x4a, 
       0x3, 0x4a, 0x3, 0x4a, 0x5, 0x4a, 0x318, 0xa, 0x4a, 0x3, 0x4b, 0x3, 
       0x4b, 0x3, 0x4b, 0x3, 0x4b, 0x3, 0x4b, 0x5, 0x4b, 0x31f, 0xa, 0x4b, 
       0x3, 0x4c, 0x3, 0x4c, 0x3, 0x4c, 0x5, 0x4c, 0x324, 0xa, 0x4c, 0x3, 
       0x4d, 0x3, 0x4d, 0x3, 0x4d, 0x3, 0x4d, 0x3, 0x4d, 0x3, 0x4d, 0x3, 
       0x4d, 0x3, 0x4d, 0x3, 0x4d, 0x7, 0x4d, 0x32f, 0xa, 0x4d, 0xc, 0x4d, 
       0xe, 0x4d, 0x332, 0xb, 0x4d, 0x3, 0x4e, 0x3, 0x4e, 0x3, 0x4e, 0x3, 
       0x4e, 0x3, 0x4e, 0x3, 0x4e, 0x3, 0x4e, 0x3, 0x4e, 0x3, 0x4e, 0x3, 
       0x4e, 0x5, 0x4e, 0x33e, 0xa, 0x4e, 0x3, 0x4f, 0x3, 0x4f, 0x3, 0x4f, 
       0x3, 0x4f, 0x3, 0x4f, 0x3, 0x4f, 0x3, 0x50, 0x3, 0x50, 0x3, 0x50, 
       0x3, 0x50, 0x5, 0x50, 0x34a, 0xa, 0x50, 0x5, 0x50, 0x34c, 0xa, 0x50, 
       0x3, 0x50, 0x3, 0x50, 0x3, 0x51, 0x3, 0x51, 0x3, 0x51, 0x7, 0x51, 
       0x353, 0xa, 0x51, 0xc, 0x51, 0xe, 0x51, 0x356, 0xb, 0x51, 0x3, 0x52, 
       0x3, 0x52, 0x3, 0x52, 0x3, 0x52, 0x3, 0x52, 0x3, 0x52, 0x3, 0x52, 
       0x5, 0x52, 0x35f, 0xa, 0x52, 0x3, 0x53, 0x3, 0x53, 0x3, 0x53, 0x3, 
       0x53, 0x3, 0x53, 0x3, 0x53, 0x5, 0x53, 0x367, 0xa, 0x53, 0x3, 0x54, 
       0x3, 0x54, 0x3, 0x54, 0x6, 0x54, 0x36c, 0xa, 0x54, 0xd, 0x54, 0xe, 
       0x54, 0x36d, 0x3, 0x55, 0x3, 0x55, 0x3, 0x56, 0x3, 0x56, 0x3, 0x56, 
       0x3, 0x56, 0x5, 0x56, 0x376, 0xa, 0x56, 0x3, 0x57, 0x3, 0x57, 0x3, 
       0x57, 0x7, 0x57, 0x37b, 0xa, 0x57, 0xc, 0x57, 0xe, 0x57, 0x37e, 0xb, 
       0x57, 0x3, 0x57, 0x2, 0xb, 0x44, 0x80, 0x82, 0x84, 0x86, 0x88, 0x8a, 
       0x8e, 0x98, 0x58, 0x2, 0x4, 0x6, 0x8, 0xa, 0xc, 0xe, 0x10, 0x12, 
       0x14, 0x16, 0x18, 0x1a, 0x1c, 0x1e, 0x20, 0x22, 0x24, 0x26, 0x28, 
       0x2a, 0x2c, 0x2e, 0x30, 0x32, 0x34, 0x36, 0x38, 0x3a, 0x3c, 0x3e, 
       0x40, 0x42, 0x44, 0x46, 0x48, 0x4a, 0x4c, 0x4e, 0x50, 0x52, 0x54, 
       0x56, 0x58, 0x5a, 0x5c, 0x5e, 0x60, 0x62, 0x64, 0x66, 0x68, 0x6a, 
       0x6c, 0x6e, 0x70, 0x72, 0x74, 0x76, 0x78, 0x7a, 0x7c, 0x7e, 0x80, 
       0x82, 0x84, 0x86, 0x88, 0x8a, 0x8c, 0x8e, 0x90, 0x92, 0x94, 0x96, 
       0x98, 0x9a, 0x9c, 0x9e, 0xa0, 0xa2, 0xa4, 0xa6, 0xa8, 0xaa, 0xac, 
       0x2, 0x9, 0x3, 0x2, 0x22, 0x23, 0x3, 0x2, 0x24, 0x27, 0x3, 0x2, 0x28, 
       0x29, 0x3, 0x2, 0x2a, 0x2b, 0x4, 0x2, 0x29, 0x29, 0x2c, 0x2c, 0x3, 
       0x2, 0x20, 0x2c, 0x4, 0x2, 0x14, 0x14, 0x3a, 0x3a, 0x2, 0x3a0, 0x2, 
       0xcc, 0x3, 0x2, 0x2, 0x2, 0x4, 0xce, 0x3, 0x2, 0x2, 0x2, 0x6, 0xd1, 
       0x3, 0x2, 0x2, 0x2, 0x8, 0xd4, 0x3, 0x2, 0x2, 0x2, 0xa, 0xd9, 0x3, 
       0x2, 0x2, 0x2, 0xc, 0xec, 0x3, 0x2, 0x2, 0x2, 0xe, 0xee, 0x3, 0x2, 
       0x2, 0x2, 0x10, 0xf7, 0x3, 0x2, 0x2, 0x2, 0x12, 0x100, 0x3, 0x2, 
       0x2, 0x2, 0x14, 0x109, 0x3, 0x2, 0x2, 0x2, 0x16, 0x111, 0x3, 0x2, 
       0x2, 0x2, 0x18, 0x114, 0x3, 0x2, 0x2, 0x2, 0x1a, 0x120, 0x3, 0x2, 
       0x2, 0x2, 0x1c, 0x122, 0x3, 0x2, 0x2, 0x2, 0x1e, 0x130, 0x3, 0x2, 
       0x2, 0x2, 0x20, 0x154, 0x3, 0x2, 0x2, 0x2, 0x22, 0x15b, 0x3, 0x2, 
       0x2, 0x2, 0x24, 0x15d, 0x3, 0x2, 0x2, 0x2, 0x26, 0x167, 0x3, 0x2, 
       0x2, 0x2, 0x28, 0x174, 0x3, 0x2, 0x2, 0x2, 0x2a, 0x17b, 0x3, 0x2, 
       0x2, 0x2, 0x2c, 0x17f, 0x3, 0x2, 0x2, 0x2, 0x2e, 0x184, 0x3, 0x2, 
       0x2, 0x2, 0x30, 0x186, 0x3, 0x2, 0x2, 0x2, 0x32, 0x197, 0x3, 0x2, 
       0x2, 0x2, 0x34, 0x199, 0x3, 0x2, 0x2, 0x2, 0x36, 0x1a4, 0x3, 0x2, 
       0x2, 0x2, 0x38, 0x1a9, 0x3, 0x2, 0x2, 0x2, 0x3a, 0x1ac, 0x3, 0x2, 
       0x2, 0x2, 0x3c, 0x1b4, 0x3, 0x2, 0x2, 0x2, 0x3e, 0x1bb, 0x3, 0x2, 
       0x2, 0x2, 0x40, 0x1bd, 0x3, 0x2, 0x2, 0x2, 0x42, 0x1c5, 0x3, 0x2, 
       0x2, 0x2, 0x44, 0x1c7, 0x3, 0x2, 0x2, 0x2, 0x46, 0x1dc, 0x3, 0x2, 
       0x2, 0x2, 0x48, 0x1ef, 0x3, 0x2, 0x2, 0x2, 0x4a, 0x1f1, 0x3, 0x2, 
       0x2, 0x2, 0x4c, 0x1fb, 0x3, 0x2, 0x2, 0x2, 0x4e, 0x203, 0x3, 0x2, 
       0x2, 0x2, 0x50, 0x20d, 0x3, 0x2, 0x2, 0x2, 0x52, 0x20f, 0x3, 0x2, 
       0x2, 0x2, 0x54, 0x216, 0x3, 0x2, 0x2, 0x2, 0x56, 0x220, 0x3, 0x2, 
       0x2, 0x2, 0x58, 0x228, 0x3, 0x2, 0x2, 0x2, 0x5a, 0x22d, 0x3, 0x2, 
       0x2, 0x2, 0x5c, 0x232, 0x3, 0x2, 0x2, 0x2, 0x5e, 0x23a, 0x3, 0x2, 
       0x2, 0x2, 0x60, 0x23e, 0x3, 0x2, 0x2, 0x2, 0x62, 0x243, 0x3, 0x2, 
       0x2, 0x2, 0x64, 0x246, 0x3, 0x2, 0x2, 0x2, 0x66, 0x250, 0x3, 0x2, 
       0x2, 0x2, 0x68, 0x254, 0x3, 0x2, 0x2, 0x2, 0x6a, 0x25b, 0x3, 0x2, 
       0x2, 0x2, 0x6c, 0x264, 0x3, 0x2, 0x2, 0x2, 0x6e, 0x26c, 0x3, 0x2, 
       0x2, 0x2, 0x70, 0x271, 0x3, 0x2, 0x2, 0x2, 0x72, 0x27a, 0x3, 0x2, 
       0x2, 0x2, 0x74, 0x283, 0x3, 0x2, 0x2, 0x2, 0x76, 0x286, 0x3, 0x2, 
       0x2, 0x2, 0x78, 0x28c, 0x3, 0x2, 0x2, 0x2, 0x7a, 0x29c, 0x3, 0x2, 
       0x2, 0x2, 0x7c, 0x2a5, 0x3, 0x2, 0x2, 0x2, 0x7e, 0x2a7, 0x3, 0x2, 
       0x2, 0x2, 0x80, 0x2a9, 0x3, 0x2, 0x2, 0x2, 0x82, 0x2b4, 0x3, 0x2, 
       0x2, 0x2, 0x84, 0x2bf, 0x3, 0x2, 0x2, 0x2, 0x86, 0x2cd, 0x3, 0x2, 
       0x2, 0x2, 0x88, 0x2db, 0x3, 0x2, 0x2, 0x2, 0x8a, 0x2e9, 0x3, 0x2, 
       0x2, 0x2, 0x8c, 0x2fd, 0x3, 0x2, 0x2, 0x2, 0x8e, 0x2ff, 0x3, 0x2, 
       0x2, 0x2, 0x90, 0x30d, 0x3, 0x2, 0x2, 0x2, 0x92, 0x317, 0x3, 0x2, 
       0x2, 0x2, 0x94, 0x31e, 0x3, 0x2, 0x2, 0x2, 0x96, 0x323, 0x3, 0x2, 
       0x2, 0x2, 0x98, 0x325, 0x3, 0x2, 0x2, 0x2, 0x9a, 0x33d, 0x3, 0x2, 
       0x2, 0x2, 0x9c, 0x33f, 0x3, 0x2, 0x2, 0x2, 0x9e, 0x345, 0x3, 0x2, 
       0x2, 0x2, 0xa0, 0x34f, 0x3, 0x2, 0x2, 0x2, 0xa2, 0x35e, 0x3, 0x2, 
       0x2, 0x2, 0xa4, 0x366, 0x3, 0x2, 0x2, 0x2, 0xa6, 0x368, 0x3, 0x2, 
       0x2, 0x2, 0xa8, 0x36f, 0x3, 0x2, 0x2, 0x2, 0xaa, 0x375, 0x3, 0x2, 
       0x2, 0x2, 0xac, 0x377, 0x3, 0x2, 0x2, 0x2, 0xae, 0xb0, 0x5, 0x4, 
       0x3, 0x2, 0xaf, 0xb1, 0x5, 0xe, 0x8, 0x2, 0xb0, 0xaf, 0x3, 0x2, 0x2, 
       0x2, 0xb0, 0xb1, 0x3, 0x2, 0x2, 0x2, 0xb1, 0xb3, 0x3, 0x2, 0x2, 0x2, 
       0xb2, 0xb4, 0x5, 0x10, 0x9, 0x2, 0xb3, 0xb2, 0x3, 0x2, 0x2, 0x2, 
       0xb3, 0xb4, 0x3, 0x2, 0x2, 0x2, 0xb4, 0xb6, 0x3, 0x2, 0x2, 0x2, 0xb5, 
       0xb7, 0x5, 0x12, 0xa, 0x2, 0xb6, 0xb5, 0x3, 0x2, 0x2, 0x2, 0xb6, 
       0xb7, 0x3, 0x2, 0x2, 0x2, 0xb7, 0xb9, 0x3, 0x2, 0x2, 0x2, 0xb8, 0xba, 
       0x5, 0x18, 0xd, 0x2, 0xb9, 0xb8, 0x3, 0x2, 0x2, 0x2, 0xb9, 0xba, 
       0x3, 0x2, 0x2, 0x2, 0xba, 0xbb, 0x3, 0x2, 0x2, 0x2, 0xbb, 0xbc, 0x7, 
       0x2, 0x2, 0x3, 0xbc, 0xcd, 0x3, 0x2, 0x2, 0x2, 0xbd, 0xbf, 0x5, 0x6, 
       0x4, 0x2, 0xbe, 0xc0, 0x5, 0xe, 0x8, 0x2, 0xbf, 0xbe, 0x3, 0x2, 0x2, 
       0x2, 0xbf, 0xc0, 0x3, 0x2, 0x2, 0x2, 0xc0, 0xc2, 0x3, 0x2, 0x2, 0x2, 
       0xc1, 0xc3, 0x5, 0x10, 0x9, 0x2, 0xc2, 0xc1, 0x3, 0x2, 0x2, 0x2, 
       0xc2, 0xc3, 0x3, 0x2, 0x2, 0x2, 0xc3, 0xc5, 0x3, 0x2, 0x2, 0x2, 0xc4, 
       0xc6, 0x5, 0x12, 0xa, 0x2, 0xc5, 0xc4, 0x3, 0x2, 0x2, 0x2, 0xc5, 
       0xc6, 0x3, 0x2, 0x2, 0x2, 0xc6, 0xc8, 0x3, 0x2, 0x2, 0x2, 0xc7, 0xc9, 
       0x5, 0x18, 0xd, 0x2, 0xc8, 0xc7, 0x3, 0x2, 0x2, 0x2, 0xc8, 0xc9, 
       0x3, 0x2, 0x2, 0x2, 0xc9, 0xca, 0x3, 0x2, 0x2, 0x2, 0xca, 0xcb, 0x7, 
       0x2, 0x2, 0x3, 0xcb, 0xcd, 0x3, 0x2, 0x2, 0x2, 0xcc, 0xae, 0x3, 0x2, 
       0x2, 0x2, 0xcc, 0xbd, 0x3, 0x2, 0x2, 0x2, 0xcd, 0x3, 0x3, 0x2, 0x2, 
       0x2, 0xce, 0xcf, 0x7, 0x6, 0x2, 0x2, 0xcf, 0xd0, 0x5, 0x14, 0xb, 
       0x2, 0xd0, 0x5, 0x3, 0x2, 0x2, 0x2, 0xd1, 0xd2, 0x7, 0x7, 0x2, 0x2, 
       0xd2, 0xd3, 0x5, 0x14, 0xb, 0x2, 0xd3, 0x7, 0x3, 0x2, 0x2, 0x2, 0xd4, 
       0xd5, 0x7, 0x3, 0x2, 0x2, 0xd5, 0xd6, 0x7, 0x33, 0x2, 0x2, 0xd6, 
       0xd7, 0x5, 0xa, 0x6, 0x2, 0xd7, 0xd8, 0x7, 0x34, 0x2, 0x2, 0xd8, 
       0x9, 0x3, 0x2, 0x2, 0x2, 0xd9, 0xda, 0x7, 0x4, 0x2, 0x2, 0xda, 0xdc, 
       0x7, 0x33, 0x2, 0x2, 0xdb, 0xdd, 0x5, 0xc, 0x7, 0x2, 0xdc, 0xdb, 
       0x3, 0x2, 0x2, 0x2, 0xdd, 0xde, 0x3, 0x2, 0x2, 0x2, 0xde, 0xdc, 0x3, 
       0x2, 0x2, 0x2, 0xde, 0xdf, 0x3, 0x2, 0x2, 0x2, 0xdf, 0xe0, 0x3, 0x2, 
       0x2, 0x2, 0xe0, 0xe1, 0x7, 0x34, 0x2, 0x2, 0xe1, 0xb, 0x3, 0x2, 0x2, 
       0x2, 0xe2, 0xe3, 0x7, 0x5, 0x2, 0x2, 0xe3, 0xe4, 0x5, 0x14, 0xb, 
       0x2, 0xe4, 0xe5, 0x7, 0xb, 0x2, 0x2, 0xe5, 0xe6, 0x5, 0x16, 0xc, 
       0x2, 0xe6, 0xed, 0x3, 0x2, 0x2, 0x2, 0xe7, 0xe8, 0x7, 0x6, 0x2, 0x2, 
       0xe8, 0xe9, 0x5, 0x14, 0xb, 0x2, 0xe9, 0xea, 0x7, 0xb, 0x2, 0x2, 
       0xea, 0xeb, 0x5, 0x16, 0xc, 0x2, 0xeb, 0xed, 0x3, 0x2, 0x2, 0x2, 
       0xec, 0xe2, 0x3, 0x2, 0x2, 0x2, 0xec, 0xe7, 0x3, 0x2, 0x2, 0x2, 0xed, 
       0xd, 0x3, 0x2, 0x2, 0x2, 0xee, 0xef, 0x7, 0x8, 0x2, 0x2, 0xef, 0xf1, 
       0x7, 0x33, 0x2, 0x2, 0xf0, 0xf2, 0x5, 0x14, 0xb, 0x2, 0xf1, 0xf0, 
       0x3, 0x2, 0x2, 0x2, 0xf2, 0xf3, 0x3, 0x2, 0x2, 0x2, 0xf3, 0xf1, 0x3, 
       0x2, 0x2, 0x2, 0xf3, 0xf4, 0x3, 0x2, 0x2, 0x2, 0xf4, 0xf5, 0x3, 0x2, 
       0x2, 0x2, 0xf5, 0xf6, 0x7, 0x34, 0x2, 0x2, 0xf6, 0xf, 0x3, 0x2, 0x2, 
       0x2, 0xf7, 0xf8, 0x7, 0x9, 0x2, 0x2, 0xf8, 0xfa, 0x7, 0x33, 0x2, 
       0x2, 0xf9, 0xfb, 0x5, 0x16, 0xc, 0x2, 0xfa, 0xf9, 0x3, 0x2, 0x2, 
       0x2, 0xfb, 0xfc, 0x3, 0x2, 0x2, 0x2, 0xfc, 0xfa, 0x3, 0x2, 0x2, 0x2, 
       0xfc, 0xfd, 0x3, 0x2, 0x2, 0x2, 0xfd, 0xfe, 0x3, 0x2, 0x2, 0x2, 0xfe, 
       0xff, 0x7, 0x34, 0x2, 0x2, 0xff, 0x11, 0x3, 0x2, 0x2, 0x2, 0x100, 
       0x101, 0x7, 0xa, 0x2, 0x2, 0x101, 0x103, 0x7, 0x33, 0x2, 0x2, 0x102, 
       0x104, 0x5, 0x16, 0xc, 0x2, 0x103, 0x102, 0x3, 0x2, 0x2, 0x2, 0x104, 
       0x105, 0x3, 0x2, 0x2, 0x2, 0x105, 0x103, 0x3, 0x2, 0x2, 0x2, 0x105, 
       0x106, 0x3, 0x2, 0x2, 0x2, 0x106, 0x107, 0x3, 0x2, 0x2, 0x2, 0x107, 
       0x108, 0x7, 0x34, 0x2, 0x2, 0x108, 0x13, 0x3, 0x2, 0x2, 0x2, 0x109, 
       0x10e, 0x7, 0x3a, 0x2, 0x2, 0x10a, 0x10b, 0x7, 0x2e, 0x2, 0x2, 0x10b, 
       0x10d, 0x7, 0x3a, 0x2, 0x2, 0x10c, 0x10a, 0x3, 0x2, 0x2, 0x2, 0x10d, 
       0x110, 0x3, 0x2, 0x2, 0x2, 0x10e, 0x10c, 0x3, 0x2, 0x2, 0x2, 0x10e, 
       0x10f, 0x3, 0x2, 0x2, 0x2, 0x10f, 0x15, 0x3, 0x2, 0x2, 0x2, 0x110, 
       0x10e, 0x3, 0x2, 0x2, 0x2, 0x111, 0x112, 0x7, 0x3e, 0x2, 0x2, 0x112, 
       0x17, 0x3, 0x2, 0x2, 0x2, 0x113, 0x115, 0x5, 0x1a, 0xe, 0x2, 0x114, 
       0x113, 0x3, 0x2, 0x2, 0x2, 0x115, 0x116, 0x3, 0x2, 0x2, 0x2, 0x116, 
       0x114, 0x3, 0x2, 0x2, 0x2, 0x116, 0x117, 0x3, 0x2, 0x2, 0x2, 0x117, 
       0x19, 0x3, 0x2, 0x2, 0x2, 0x118, 0x121, 0x5, 0x22, 0x12, 0x2, 0x119, 
       0x121, 0x5, 0x24, 0x13, 0x2, 0x11a, 0x11b, 0x7, 0xc, 0x2, 0x2, 0x11b, 
       0x121, 0x5, 0x24, 0x13, 0x2, 0x11c, 0x121, 0x5, 0x1c, 0xf, 0x2, 0x11d, 
       0x11e, 0x7, 0xc, 0x2, 0x2, 0x11e, 0x121, 0x5, 0x1c, 0xf, 0x2, 0x11f, 
       0x121, 0x5, 0x20, 0x11, 0x2, 0x120, 0x118, 0x3, 0x2, 0x2, 0x2, 0x120, 
       0x119, 0x3, 0x2, 0x2, 0x2, 0x120, 0x11a, 0x3, 0x2, 0x2, 0x2, 0x120, 
       0x11c, 0x3, 0x2, 0x2, 0x2, 0x120, 0x11d, 0x3, 0x2, 0x2, 0x2, 0x120, 
       0x11f, 0x3, 0x2, 0x2, 0x2, 0x121, 0x1b, 0x3, 0x2, 0x2, 0x2, 0x122, 
       0x123, 0x7, 0x1a, 0x2, 0x2, 0x123, 0x125, 0x5, 0xaa, 0x56, 0x2, 0x124, 
       0x126, 0x5, 0x34, 0x1b, 0x2, 0x125, 0x124, 0x3, 0x2, 0x2, 0x2, 0x125, 
       0x126, 0x3, 0x2, 0x2, 0x2, 0x126, 0x127, 0x3, 0x2, 0x2, 0x2, 0x127, 
       0x12b, 0x7, 0x33, 0x2, 0x2, 0x128, 0x12a, 0x5, 0x1e, 0x10, 0x2, 0x129, 
       0x128, 0x3, 0x2, 0x2, 0x2, 0x12a, 0x12d, 0x3, 0x2, 0x2, 0x2, 0x12b, 
       0x129, 0x3, 0x2, 0x2, 0x2, 0x12b, 0x12c, 0x3, 0x2, 0x2, 0x2, 0x12c, 
       0x12e, 0x3, 0x2, 0x2, 0x2, 0x12d, 0x12b, 0x3, 0x2, 0x2, 0x2, 0x12e, 
       0x12f, 0x7, 0x34, 0x2, 0x2, 0x12f, 0x1d, 0x3, 0x2, 0x2, 0x2, 0x130, 
       0x131, 0x7, 0xe, 0x2, 0x2, 0x131, 0x132, 0x5, 0xaa, 0x56, 0x2, 0x132, 
       0x133, 0x5, 0x26, 0x14, 0x2, 0x133, 0x134, 0x7, 0x1d, 0x2, 0x2, 0x134, 
       0x135, 0x5, 0x3c, 0x1f, 0x2, 0x135, 0x1f, 0x3, 0x2, 0x2, 0x2, 0x136, 
       0x138, 0x7, 0x1b, 0x2, 0x2, 0x137, 0x139, 0x5, 0x34, 0x1b, 0x2, 0x138, 
       0x137, 0x3, 0x2, 0x2, 0x2, 0x138, 0x139, 0x3, 0x2, 0x2, 0x2, 0x139, 
       0x13a, 0x3, 0x2, 0x2, 0x2, 0x13a, 0x13b, 0x5, 0x3c, 0x1f, 0x2, 0x13b, 
       0x13c, 0x7, 0x19, 0x2, 0x2, 0x13c, 0x13d, 0x5, 0x3c, 0x1f, 0x2, 0x13d, 
       0x141, 0x7, 0x33, 0x2, 0x2, 0x13e, 0x140, 0x5, 0x24, 0x13, 0x2, 0x13f, 
       0x13e, 0x3, 0x2, 0x2, 0x2, 0x140, 0x143, 0x3, 0x2, 0x2, 0x2, 0x141, 
       0x13f, 0x3, 0x2, 0x2, 0x2, 0x141, 0x142, 0x3, 0x2, 0x2, 0x2, 0x142, 
       0x144, 0x3, 0x2, 0x2, 0x2, 0x143, 0x141, 0x3, 0x2, 0x2, 0x2, 0x144, 
       0x145, 0x7, 0x34, 0x2, 0x2, 0x145, 0x155, 0x3, 0x2, 0x2, 0x2, 0x146, 
       0x148, 0x7, 0x1b, 0x2, 0x2, 0x147, 0x149, 0x5, 0x34, 0x1b, 0x2, 0x148, 
       0x147, 0x3, 0x2, 0x2, 0x2, 0x148, 0x149, 0x3, 0x2, 0x2, 0x2, 0x149, 
       0x14a, 0x3, 0x2, 0x2, 0x2, 0x14a, 0x14b, 0x5, 0x3c, 0x1f, 0x2, 0x14b, 
       0x14f, 0x7, 0x33, 0x2, 0x2, 0x14c, 0x14e, 0x5, 0x24, 0x13, 0x2, 0x14d, 
       0x14c, 0x3, 0x2, 0x2, 0x2, 0x14e, 0x151, 0x3, 0x2, 0x2, 0x2, 0x14f, 
       0x14d, 0x3, 0x2, 0x2, 0x2, 0x14f, 0x150, 0x3, 0x2, 0x2, 0x2, 0x150, 
       0x152, 0x3, 0x2, 0x2, 0x2, 0x151, 0x14f, 0x3, 0x2, 0x2, 0x2, 0x152, 
       0x153, 0x7, 0x34, 0x2, 0x2, 0x153, 0x155, 0x3, 0x2, 0x2, 0x2, 0x154, 
       0x136, 0x3, 0x2, 0x2, 0x2, 0x154, 0x146, 0x3, 0x2, 0x2, 0x2, 0x155, 
       0x21, 0x3, 0x2, 0x2, 0x2, 0x156, 0x15c, 0x5, 0x2e, 0x18, 0x2, 0x157, 
       0x158, 0x7, 0xd, 0x2, 0x2, 0x158, 0x15c, 0x5, 0x30, 0x19, 0x2, 0x159, 
       0x15a, 0x7, 0xc, 0x2, 0x2, 0x15a, 0x15c, 0x5, 0x30, 0x19, 0x2, 0x15b, 
       0x156, 0x3, 0x2, 0x2, 0x2, 0x15b, 0x157, 0x3, 0x2, 0x2, 0x2, 0x15b, 
       0x159, 0x3, 0x2, 0x2, 0x2, 0x15c, 0x23, 0x3, 0x2, 0x2, 0x2, 0x15d, 
       0x15e, 0x7, 0xe, 0x2, 0x2, 0x15e, 0x160, 0x5, 0xaa, 0x56, 0x2, 0x15f, 
       0x161, 0x5, 0x34, 0x1b, 0x2, 0x160, 0x15f, 0x3, 0x2, 0x2, 0x2, 0x160, 
       0x161, 0x3, 0x2, 0x2, 0x2, 0x161, 0x162, 0x3, 0x2, 0x2, 0x2, 0x162, 
       0x163, 0x5, 0x26, 0x14, 0x2, 0x163, 0x164, 0x7, 0x1d, 0x2, 0x2, 0x164, 
       0x165, 0x5, 0x3c, 0x1f, 0x2, 0x165, 0x166, 0x5, 0x72, 0x3a, 0x2, 
       0x166, 0x25, 0x3, 0x2, 0x2, 0x2, 0x167, 0x170, 0x7, 0x35, 0x2, 0x2, 
       0x168, 0x16d, 0x5, 0x28, 0x15, 0x2, 0x169, 0x16a, 0x7, 0x30, 0x2, 
       0x2, 0x16a, 0x16c, 0x5, 0x28, 0x15, 0x2, 0x16b, 0x169, 0x3, 0x2, 
       0x2, 0x2, 0x16c, 0x16f, 0x3, 0x2, 0x2, 0x2, 0x16d, 0x16b, 0x3, 0x2, 
       0x2, 0x2, 0x16d, 0x16e, 0x3, 0x2, 0x2, 0x2, 0x16e, 0x171, 0x3, 0x2, 
       0x2, 0x2, 0x16f, 0x16d, 0x3, 0x2, 0x2, 0x2, 0x170, 0x168, 0x3, 0x2, 
       0x2, 0x2, 0x170, 0x171, 0x3, 0x2, 0x2, 0x2, 0x171, 0x172, 0x3, 0x2, 
       0x2, 0x2, 0x172, 0x173, 0x7, 0x36, 0x2, 0x2, 0x173, 0x27, 0x3, 0x2, 
       0x2, 0x2, 0x174, 0x175, 0x7, 0x3a, 0x2, 0x2, 0x175, 0x176, 0x7, 0x32, 
       0x2, 0x2, 0x176, 0x177, 0x5, 0x3c, 0x1f, 0x2, 0x177, 0x29, 0x3, 0x2, 
       0x2, 0x2, 0x178, 0x179, 0x7, 0xc, 0x2, 0x2, 0x179, 0x17c, 0x5, 0x2c, 
       0x17, 0x2, 0x17a, 0x17c, 0x5, 0x2c, 0x17, 0x2, 0x17b, 0x178, 0x3, 
       0x2, 0x2, 0x2, 0x17b, 0x17a, 0x3, 0x2, 0x2, 0x2, 0x17c, 0x2b, 0x3, 
       0x2, 0x2, 0x2, 0x17d, 0x180, 0x5, 0x30, 0x19, 0x2, 0x17e, 0x180, 
       0x5, 0x32, 0x1a, 0x2, 0x17f, 0x17d, 0x3, 0x2, 0x2, 0x2, 0x17f, 0x17e, 
       0x3, 0x2, 0x2, 0x2, 0x180, 0x2d, 0x3, 0x2, 0x2, 0x2, 0x181, 0x182, 
       0x7, 0xc, 0x2, 0x2, 0x182, 0x185, 0x5, 0x32, 0x1a, 0x2, 0x183, 0x185, 
       0x5, 0x32, 0x1a, 0x2, 0x184, 0x181, 0x3, 0x2, 0x2, 0x2, 0x184, 0x183, 
       0x3, 0x2, 0x2, 0x2, 0x185, 0x2f, 0x3, 0x2, 0x2, 0x2, 0x186, 0x187, 
       0x5, 0xaa, 0x56, 0x2, 0x187, 0x188, 0x7, 0x32, 0x2, 0x2, 0x188, 0x189, 
       0x5, 0x3c, 0x1f, 0x2, 0x189, 0x31, 0x3, 0x2, 0x2, 0x2, 0x18a, 0x18b, 
       0x7, 0xf, 0x2, 0x2, 0x18b, 0x18d, 0x5, 0xaa, 0x56, 0x2, 0x18c, 0x18e, 
       0x5, 0x34, 0x1b, 0x2, 0x18d, 0x18c, 0x3, 0x2, 0x2, 0x2, 0x18d, 0x18e, 
       0x3, 0x2, 0x2, 0x2, 0x18e, 0x18f, 0x3, 0x2, 0x2, 0x2, 0x18f, 0x190, 
       0x7, 0x2f, 0x2, 0x2, 0x190, 0x191, 0x5, 0x3c, 0x1f, 0x2, 0x191, 0x198, 
       0x3, 0x2, 0x2, 0x2, 0x192, 0x193, 0x7, 0xf, 0x2, 0x2, 0x193, 0x195, 
       0x5, 0xaa, 0x56, 0x2, 0x194, 0x196, 0x5, 0x34, 0x1b, 0x2, 0x195, 
       0x194, 0x3, 0x2, 0x2, 0x2, 0x195, 0x196, 0x3, 0x2, 0x2, 0x2, 0x196, 
       0x198, 0x3, 0x2, 0x2, 0x2, 0x197, 0x18a, 0x3, 0x2, 0x2, 0x2, 0x197, 
       0x192, 0x3, 0x2, 0x2, 0x2, 0x198, 0x33, 0x3, 0x2, 0x2, 0x2, 0x199, 
       0x19a, 0x7, 0x37, 0x2, 0x2, 0x19a, 0x19f, 0x5, 0x36, 0x1c, 0x2, 0x19b, 
       0x19c, 0x7, 0x30, 0x2, 0x2, 0x19c, 0x19e, 0x5, 0x36, 0x1c, 0x2, 0x19d, 
       0x19b, 0x3, 0x2, 0x2, 0x2, 0x19e, 0x1a1, 0x3, 0x2, 0x2, 0x2, 0x19f, 
       0x19d, 0x3, 0x2, 0x2, 0x2, 0x19f, 0x1a0, 0x3, 0x2, 0x2, 0x2, 0x1a0, 
       0x1a2, 0x3, 0x2, 0x2, 0x2, 0x1a1, 0x19f, 0x3, 0x2, 0x2, 0x2, 0x1a2, 
       0x1a3, 0x7, 0x38, 0x2, 0x2, 0x1a3, 0x35, 0x3, 0x2, 0x2, 0x2, 0x1a4, 
       0x1a7, 0x5, 0x38, 0x1d, 0x2, 0x1a5, 0x1a6, 0x7, 0x32, 0x2, 0x2, 0x1a6, 
       0x1a8, 0x5, 0x3a, 0x1e, 0x2, 0x1a7, 0x1a5, 0x3, 0x2, 0x2, 0x2, 0x1a7, 
       0x1a8, 0x3, 0x2, 0x2, 0x2, 0x1a8, 0x37, 0x3, 0x2, 0x2, 0x2, 0x1a9, 
       0x1aa, 0x7, 0x39, 0x2, 0x2, 0x1aa, 0x1ab, 0x7, 0x3a, 0x2, 0x2, 0x1ab, 
       0x39, 0x3, 0x2, 0x2, 0x2, 0x1ac, 0x1b1, 0x5, 0x3c, 0x1f, 0x2, 0x1ad, 
       0x1ae, 0x7, 0x28, 0x2, 0x2, 0x1ae, 0x1b0, 0x5, 0x3c, 0x1f, 0x2, 0x1af, 
       0x1ad, 0x3, 0x2, 0x2, 0x2, 0x1b0, 0x1b3, 0x3, 0x2, 0x2, 0x2, 0x1b1, 
       0x1af, 0x3, 0x2, 0x2, 0x2, 0x1b1, 0x1b2, 0x3, 0x2, 0x2, 0x2, 0x1b2, 
       0x3b, 0x3, 0x2, 0x2, 0x2, 0x1b3, 0x1b1, 0x3, 0x2, 0x2, 0x2, 0x1b4, 
       0x1b5, 0x5, 0x3e, 0x20, 0x2, 0x1b5, 0x3d, 0x3, 0x2, 0x2, 0x2, 0x1b6, 
       0x1b7, 0x7, 0x10, 0x2, 0x2, 0x1b7, 0x1b8, 0x5, 0x34, 0x1b, 0x2, 0x1b8, 
       0x1b9, 0x5, 0x40, 0x21, 0x2, 0x1b9, 0x1bc, 0x3, 0x2, 0x2, 0x2, 0x1ba, 
       0x1bc, 0x5, 0x40, 0x21, 0x2, 0x1bb, 0x1b6, 0x3, 0x2, 0x2, 0x2, 0x1bb, 
       0x1ba, 0x3, 0x2, 0x2, 0x2, 0x1bc, 0x3f, 0x3, 0x2, 0x2, 0x2, 0x1bd, 
       0x1c0, 0x5, 0x42, 0x22, 0x2, 0x1be, 0x1bf, 0x7, 0x1d, 0x2, 0x2, 0x1bf, 
       0x1c1, 0x5, 0x40, 0x21, 0x2, 0x1c0, 0x1be, 0x3, 0x2, 0x2, 0x2, 0x1c0, 
       0x1c1, 0x3, 0x2, 0x2, 0x2, 0x1c1, 0x41, 0x3, 0x2, 0x2, 0x2, 0x1c2, 
       0x1c3, 0x7, 0x2d, 0x2, 0x2, 0x1c3, 0x1c6, 0x5, 0x42, 0x22, 0x2, 0x1c4, 
       0x1c6, 0x5, 0x44, 0x23, 0x2, 0x1c5, 0x1c2, 0x3, 0x2, 0x2, 0x2, 0x1c5, 
       0x1c4, 0x3, 0x2, 0x2, 0x2, 0x1c6, 0x43, 0x3, 0x2, 0x2, 0x2, 0x1c7, 
       0x1c8, 0x8, 0x23, 0x1, 0x2, 0x1c8, 0x1c9, 0x5, 0x46, 0x24, 0x2, 0x1c9, 
       0x1ce, 0x3, 0x2, 0x2, 0x2, 0x1ca, 0x1cb, 0xc, 0x4, 0x2, 0x2, 0x1cb, 
       0x1cd, 0x5, 0x46, 0x24, 0x2, 0x1cc, 0x1ca, 0x3, 0x2, 0x2, 0x2, 0x1cd, 
       0x1d0, 0x3, 0x2, 0x2, 0x2, 0x1ce, 0x1cc, 0x3, 0x2, 0x2, 0x2, 0x1ce, 
       0x1cf, 0x3, 0x2, 0x2, 0x2, 0x1cf, 0x45, 0x3, 0x2, 0x2, 0x2, 0x1d0, 
       0x1ce, 0x3, 0x2, 0x2, 0x2, 0x1d1, 0x1dd, 0x5, 0x48, 0x25, 0x2, 0x1d2, 
       0x1dd, 0x5, 0x4a, 0x26, 0x2, 0x1d3, 0x1dd, 0x5, 0x50, 0x29, 0x2, 
       0x1d4, 0x1dd, 0x5, 0x54, 0x2b, 0x2, 0x1d5, 0x1d6, 0x7, 0x39, 0x2, 
       0x2, 0x1d6, 0x1dd, 0x7, 0x3a, 0x2, 0x2, 0x1d7, 0x1dd, 0x5, 0xaa, 
       0x56, 0x2, 0x1d8, 0x1d9, 0x7, 0x35, 0x2, 0x2, 0x1d9, 0x1da, 0x5, 
       0x3c, 0x1f, 0x2, 0x1da, 0x1db, 0x7, 0x36, 0x2, 0x2, 0x1db, 0x1dd, 
       0x3, 0x2, 0x2, 0x2, 0x1dc, 0x1d1, 0x3, 0x2, 0x2, 0x2, 0x1dc, 0x1d2, 
       0x3, 0x2, 0x2, 0x2, 0x1dc, 0x1d3, 0x3, 0x2, 0x2, 0x2, 0x1dc, 0x1d4, 
       0x3, 0x2, 0x2, 0x2, 0x1dc, 0x1d5, 0x3, 0x2, 0x2, 0x2, 0x1dc, 0x1d7, 
       0x3, 0x2, 0x2, 0x2, 0x1dc, 0x1d8, 0x3, 0x2, 0x2, 0x2, 0x1dd, 0x47, 
       0x3, 0x2, 0x2, 0x2, 0x1de, 0x1df, 0x7, 0x37, 0x2, 0x2, 0x1df, 0x1e0, 
       0x5, 0x3c, 0x1f, 0x2, 0x1e0, 0x1e1, 0x7, 0x30, 0x2, 0x2, 0x1e1, 0x1e2, 
       0x7, 0x3b, 0x2, 0x2, 0x1e2, 0x1e3, 0x7, 0x38, 0x2, 0x2, 0x1e3, 0x1f0, 
       0x3, 0x2, 0x2, 0x2, 0x1e4, 0x1e5, 0x7, 0x37, 0x2, 0x2, 0x1e5, 0x1e6, 
       0x5, 0x3c, 0x1f, 0x2, 0x1e6, 0x1e7, 0x7, 0x38, 0x2, 0x2, 0x1e7, 0x1f0, 
       0x3, 0x2, 0x2, 0x2, 0x1e8, 0x1e9, 0x7, 0x37, 0x2, 0x2, 0x1e9, 0x1ea, 
       0x5, 0x3c, 0x1f, 0x2, 0x1ea, 0x1eb, 0x7, 0x30, 0x2, 0x2, 0x1eb, 0x1ec, 
       0x7, 0x29, 0x2, 0x2, 0x1ec, 0x1ed, 0x7, 0x3b, 0x2, 0x2, 0x1ed, 0x1ee, 
       0x7, 0x38, 0x2, 0x2, 0x1ee, 0x1f0, 0x3, 0x2, 0x2, 0x2, 0x1ef, 0x1de, 
       0x3, 0x2, 0x2, 0x2, 0x1ef, 0x1e4, 0x3, 0x2, 0x2, 0x2, 0x1ef, 0x1e8, 
       0x3, 0x2, 0x2, 0x2, 0x1f0, 0x49, 0x3, 0x2, 0x2, 0x2, 0x1f1, 0x1f2, 
       0x7, 0x11, 0x2, 0x2, 0x1f2, 0x1f7, 0x7, 0x33, 0x2, 0x2, 0x1f3, 0x1f5, 
       0x5, 0x4c, 0x27, 0x2, 0x1f4, 0x1f6, 0x7, 0x30, 0x2, 0x2, 0x1f5, 0x1f4, 
       0x3, 0x2, 0x2, 0x2, 0x1f5, 0x1f6, 0x3, 0x2, 0x2, 0x2, 0x1f6, 0x1f8, 
       0x3, 0x2, 0x2, 0x2, 0x1f7, 0x1f3, 0x3, 0x2, 0x2, 0x2, 0x1f7, 0x1f8, 
       0x3, 0x2, 0x2, 0x2, 0x1f8, 0x1f9, 0x3, 0x2, 0x2, 0x2, 0x1f9, 0x1fa, 
       0x7, 0x34, 0x2, 0x2, 0x1fa, 0x4b, 0x3, 0x2, 0x2, 0x2, 0x1fb, 0x200, 
       0x5, 0x4e, 0x28, 0x2, 0x1fc, 0x1fd, 0x7, 0x30, 0x2, 0x2, 0x1fd, 0x1ff, 
       0x5, 0x4e, 0x28, 0x2, 0x1fe, 0x1fc, 0x3, 0x2, 0x2, 0x2, 0x1ff, 0x202, 
       0x3, 0x2, 0x2, 0x2, 0x200, 0x1fe, 0x3, 0x2, 0x2, 0x2, 0x200, 0x201, 
       0x3, 0x2, 0x2, 0x2, 0x201, 0x4d, 0x3, 0x2, 0x2, 0x2, 0x202, 0x200, 
       0x3, 0x2, 0x2, 0x2, 0x203, 0x204, 0x5, 0xaa, 0x56, 0x2, 0x204, 0x205, 
       0x7, 0x32, 0x2, 0x2, 0x205, 0x206, 0x5, 0x3c, 0x1f, 0x2, 0x206, 0x4f, 
       0x3, 0x2, 0x2, 0x2, 0x207, 0x208, 0x7, 0x35, 0x2, 0x2, 0x208, 0x20e, 
       0x7, 0x36, 0x2, 0x2, 0x209, 0x20a, 0x7, 0x35, 0x2, 0x2, 0x20a, 0x20b, 
       0x5, 0x52, 0x2a, 0x2, 0x20b, 0x20c, 0x7, 0x36, 0x2, 0x2, 0x20c, 0x20e, 
       0x3, 0x2, 0x2, 0x2, 0x20d, 0x207, 0x3, 0x2, 0x2, 0x2, 0x20d, 0x209, 
       0x3, 0x2, 0x2, 0x2, 0x20e, 0x51, 0x3, 0x2, 0x2, 0x2, 0x20f, 0x212, 
       0x5, 0x3c, 0x1f, 0x2, 0x210, 0x211, 0x7, 0x30, 0x2, 0x2, 0x211, 0x213, 
       0x5, 0x3c, 0x1f, 0x2, 0x212, 0x210, 0x3, 0x2, 0x2, 0x2, 0x213, 0x214, 
       0x3, 0x2, 0x2, 0x2, 0x214, 0x212, 0x3, 0x2, 0x2, 0x2, 0x214, 0x215, 
       0x3, 0x2, 0x2, 0x2, 0x215, 0x53, 0x3, 0x2, 0x2, 0x2, 0x216, 0x217, 
       0x7, 0x12, 0x2, 0x2, 0x217, 0x21c, 0x7, 0x33, 0x2, 0x2, 0x218, 0x21a, 
       0x5, 0x56, 0x2c, 0x2, 0x219, 0x21b, 0x7, 0x30, 0x2, 0x2, 0x21a, 0x219, 
       0x3, 0x2, 0x2, 0x2, 0x21a, 0x21b, 0x3, 0x2, 0x2, 0x2, 0x21b, 0x21d, 
       0x3, 0x2, 0x2, 0x2, 0x21c, 0x218, 0x3, 0x2, 0x2, 0x2, 0x21c, 0x21d, 
       0x3, 0x2, 0x2, 0x2, 0x21d, 0x21e, 0x3, 0x2, 0x2, 0x2, 0x21e, 0x21f, 
       0x7, 0x34, 0x2, 0x2, 0x21f, 0x55, 0x3, 0x2, 0x2, 0x2, 0x220, 0x225, 
       0x5, 0x58, 0x2d, 0x2, 0x221, 0x222, 0x7, 0x30, 0x2, 0x2, 0x222, 0x224, 
       0x5, 0x58, 0x2d, 0x2, 0x223, 0x221, 0x3, 0x2, 0x2, 0x2, 0x224, 0x227, 
       0x3, 0x2, 0x2, 0x2, 0x225, 0x223, 0x3, 0x2, 0x2, 0x2, 0x225, 0x226, 
       0x3, 0x2, 0x2, 0x2, 0x226, 0x57, 0x3, 0x2, 0x2, 0x2, 0x227, 0x225, 
       0x3, 0x2, 0x2, 0x2, 0x228, 0x229, 0x5, 0xaa, 0x56, 0x2, 0x229, 0x22a, 
       0x5, 0x3c, 0x1f, 0x2, 0x22a, 0x59, 0x3, 0x2, 0x2, 0x2, 0x22b, 0x22e, 
       0x5, 0x5c, 0x2f, 0x2, 0x22c, 0x22e, 0x5, 0x5e, 0x30, 0x2, 0x22d, 
       0x22b, 0x3, 0x2, 0x2, 0x2, 0x22d, 0x22c, 0x3, 0x2, 0x2, 0x2, 0x22e, 
       0x5b, 0x3, 0x2, 0x2, 0x2, 0x22f, 0x233, 0x5, 0x60, 0x31, 0x2, 0x230, 
       0x233, 0x5, 0x7e, 0x40, 0x2, 0x231, 0x233, 0x5, 0x62, 0x32, 0x2, 
       0x232, 0x22f, 0x3, 0x2, 0x2, 0x2, 0x232, 0x230, 0x3, 0x2, 0x2, 0x2, 
       0x232, 0x231, 0x3, 0x2, 0x2, 0x2, 0x233, 0x5d, 0x3, 0x2, 0x2, 0x2, 
       0x234, 0x23b, 0x5, 0x72, 0x3a, 0x2, 0x235, 0x23b, 0x5, 0x64, 0x33, 
       0x2, 0x236, 0x23b, 0x5, 0x66, 0x34, 0x2, 0x237, 0x23b, 0x5, 0x68, 
       0x35, 0x2, 0x238, 0x23b, 0x5, 0x6a, 0x36, 0x2, 0x239, 0x23b, 0x5, 
       0x70, 0x39, 0x2, 0x23a, 0x234, 0x3, 0x2, 0x2, 0x2, 0x23a, 0x235, 
       0x3, 0x2, 0x2, 0x2, 0x23a, 0x236, 0x3, 0x2, 0x2, 0x2, 0x23a, 0x237, 
       0x3, 0x2, 0x2, 0x2, 0x23a, 0x238, 0x3, 0x2, 0x2, 0x2, 0x23a, 0x239, 
       0x3, 0x2, 0x2, 0x2, 0x23b, 0x5f, 0x3, 0x2, 0x2, 0x2, 0x23c, 0x23f, 
       0x5, 0xaa, 0x56, 0x2, 0x23d, 0x23f, 0x5, 0xa4, 0x53, 0x2, 0x23e, 
       0x23c, 0x3, 0x2, 0x2, 0x2, 0x23e, 0x23d, 0x3, 0x2, 0x2, 0x2, 0x23f, 
       0x240, 0x3, 0x2, 0x2, 0x2, 0x240, 0x241, 0x7, 0x1f, 0x2, 0x2, 0x241, 
       0x242, 0x5, 0x5a, 0x2e, 0x2, 0x242, 0x61, 0x3, 0x2, 0x2, 0x2, 0x243, 
       0x244, 0x7, 0x16, 0x2, 0x2, 0x244, 0x245, 0x5, 0x5c, 0x2f, 0x2, 0x245, 
       0x63, 0x3, 0x2, 0x2, 0x2, 0x246, 0x247, 0x7, 0x17, 0x2, 0x2, 0x247, 
       0x248, 0x5, 0x5c, 0x2f, 0x2, 0x248, 0x24e, 0x5, 0x72, 0x3a, 0x2, 
       0x249, 0x24c, 0x7, 0x18, 0x2, 0x2, 0x24a, 0x24d, 0x5, 0x72, 0x3a, 
       0x2, 0x24b, 0x24d, 0x5, 0x64, 0x33, 0x2, 0x24c, 0x24a, 0x3, 0x2, 
       0x2, 0x2, 0x24c, 0x24b, 0x3, 0x2, 0x2, 0x2, 0x24d, 0x24f, 0x3, 0x2, 
       0x2, 0x2, 0x24e, 0x249, 0x3, 0x2, 0x2, 0x2, 0x24e, 0x24f, 0x3, 0x2, 
       0x2, 0x2, 0x24f, 0x65, 0x3, 0x2, 0x2, 0x2, 0x250, 0x251, 0x7, 0x19, 
       0x2, 0x2, 0x251, 0x252, 0x5, 0x5c, 0x2f, 0x2, 0x252, 0x253, 0x5, 
       0x72, 0x3a, 0x2, 0x253, 0x67, 0x3, 0x2, 0x2, 0x2, 0x254, 0x256, 0x7, 
       0xe, 0x2, 0x2, 0x255, 0x257, 0x5, 0x34, 0x1b, 0x2, 0x256, 0x255, 
       0x3, 0x2, 0x2, 0x2, 0x256, 0x257, 0x3, 0x2, 0x2, 0x2, 0x257, 0x258, 
       0x3, 0x2, 0x2, 0x2, 0x258, 0x259, 0x5, 0x26, 0x14, 0x2, 0x259, 0x25a, 
       0x5, 0x72, 0x3a, 0x2, 0x25a, 0x69, 0x3, 0x2, 0x2, 0x2, 0x25b, 0x25c, 
       0x7, 0x13, 0x2, 0x2, 0x25c, 0x25d, 0x5, 0x5a, 0x2e, 0x2, 0x25d, 0x25e, 
       0x7, 0x33, 0x2, 0x2, 0x25e, 0x260, 0x5, 0x6c, 0x37, 0x2, 0x25f, 0x261, 
       0x7, 0x30, 0x2, 0x2, 0x260, 0x25f, 0x3, 0x2, 0x2, 0x2, 0x260, 0x261, 
       0x3, 0x2, 0x2, 0x2, 0x261, 0x262, 0x3, 0x2, 0x2, 0x2, 0x262, 0x263, 
       0x7, 0x34, 0x2, 0x2, 0x263, 0x6b, 0x3, 0x2, 0x2, 0x2, 0x264, 0x269, 
       0x5, 0x6e, 0x38, 0x2, 0x265, 0x266, 0x7, 0x30, 0x2, 0x2, 0x266, 0x268, 
       0x5, 0x6e, 0x38, 0x2, 0x267, 0x265, 0x3, 0x2, 0x2, 0x2, 0x268, 0x26b, 
       0x3, 0x2, 0x2, 0x2, 0x269, 0x267, 0x3, 0x2, 0x2, 0x2, 0x269, 0x26a, 
       0x3, 0x2, 0x2, 0x2, 0x26a, 0x6d, 0x3, 0x2, 0x2, 0x2, 0x26b, 0x269, 
       0x3, 0x2, 0x2, 0x2, 0x26c, 0x26d, 0x5, 0xaa, 0x56, 0x2, 0x26d, 0x26e, 
       0x7, 0x3a, 0x2, 0x2, 0x26e, 0x26f, 0x7, 0x1e, 0x2, 0x2, 0x26f, 0x270, 
       0x5, 0x5a, 0x2e, 0x2, 0x270, 0x6f, 0x3, 0x2, 0x2, 0x2, 0x271, 0x272, 
       0x7, 0x14, 0x2, 0x2, 0x272, 0x273, 0x5, 0x5a, 0x2e, 0x2, 0x273, 0x274, 
       0x7, 0x33, 0x2, 0x2, 0x274, 0x276, 0x5, 0xa0, 0x51, 0x2, 0x275, 0x277, 
       0x7, 0x30, 0x2, 0x2, 0x276, 0x275, 0x3, 0x2, 0x2, 0x2, 0x276, 0x277, 
       0x3, 0x2, 0x2, 0x2, 0x277, 0x278, 0x3, 0x2, 0x2, 0x2, 0x278, 0x279, 
       0x7, 0x34, 0x2, 0x2, 0x279, 0x71, 0x3, 0x2, 0x2, 0x2, 0x27a, 0x27b, 
       0x7, 0x33, 0x2, 0x2, 0x27b, 0x27c, 0x5, 0x74, 0x3b, 0x2, 0x27c, 0x27d, 
       0x7, 0x34, 0x2, 0x2, 0x27d, 0x73, 0x3, 0x2, 0x2, 0x2, 0x27e, 0x280, 
       0x5, 0x76, 0x3c, 0x2, 0x27f, 0x281, 0x5, 0x5c, 0x2f, 0x2, 0x280, 
       0x27f, 0x3, 0x2, 0x2, 0x2, 0x280, 0x281, 0x3, 0x2, 0x2, 0x2, 0x281, 
       0x284, 0x3, 0x2, 0x2, 0x2, 0x282, 0x284, 0x5, 0x5c, 0x2f, 0x2, 0x283, 
       0x27e, 0x3, 0x2, 0x2, 0x2, 0x283, 0x282, 0x3, 0x2, 0x2, 0x2, 0x284, 
       0x75, 0x3, 0x2, 0x2, 0x2, 0x285, 0x287, 0x5, 0x78, 0x3d, 0x2, 0x286, 
       0x285, 0x3, 0x2, 0x2, 0x2, 0x287, 0x288, 0x3, 0x2, 0x2, 0x2, 0x288, 
       0x286, 0x3, 0x2, 0x2, 0x2, 0x288, 0x289, 0x3, 0x2, 0x2, 0x2, 0x289, 
       0x77, 0x3, 0x2, 0x2, 0x2, 0x28a, 0x28d, 0x5, 0x7a, 0x3e, 0x2, 0x28b, 
       0x28d, 0x5, 0x7c, 0x3f, 0x2, 0x28c, 0x28a, 0x3, 0x2, 0x2, 0x2, 0x28c, 
       0x28b, 0x3, 0x2, 0x2, 0x2, 0x28d, 0x79, 0x3, 0x2, 0x2, 0x2, 0x28e, 
       0x28f, 0x7, 0x15, 0x2, 0x2, 0x28f, 0x290, 0x5, 0xaa, 0x56, 0x2, 0x290, 
       0x291, 0x7, 0x32, 0x2, 0x2, 0x291, 0x292, 0x5, 0x3c, 0x1f, 0x2, 0x292, 
       0x293, 0x7, 0x2f, 0x2, 0x2, 0x293, 0x294, 0x5, 0x5a, 0x2e, 0x2, 0x294, 
       0x295, 0x7, 0x31, 0x2, 0x2, 0x295, 0x29d, 0x3, 0x2, 0x2, 0x2, 0x296, 
       0x297, 0x7, 0x15, 0x2, 0x2, 0x297, 0x298, 0x5, 0xaa, 0x56, 0x2, 0x298, 
       0x299, 0x7, 0x2f, 0x2, 0x2, 0x299, 0x29a, 0x5, 0x5a, 0x2e, 0x2, 0x29a, 
       0x29b, 0x7, 0x31, 0x2, 0x2, 0x29b, 0x29d, 0x3, 0x2, 0x2, 0x2, 0x29c, 
       0x28e, 0x3, 0x2, 0x2, 0x2, 0x29c, 0x296, 0x3, 0x2, 0x2, 0x2, 0x29d, 
       0x7b, 0x3, 0x2, 0x2, 0x2, 0x29e, 0x29f, 0x5, 0x5c, 0x2f, 0x2, 0x29f, 
       0x2a0, 0x7, 0x31, 0x2, 0x2, 0x2a0, 0x2a6, 0x3, 0x2, 0x2, 0x2, 0x2a1, 
       0x2a3, 0x5, 0x5e, 0x30, 0x2, 0x2a2, 0x2a4, 0x7, 0x31, 0x2, 0x2, 0x2a3, 
       0x2a2, 0x3, 0x2, 0x2, 0x2, 0x2a3, 0x2a4, 0x3, 0x2, 0x2, 0x2, 0x2a4, 
       0x2a6, 0x3, 0x2, 0x2, 0x2, 0x2a5, 0x29e, 0x3, 0x2, 0x2, 0x2, 0x2a5, 
       0x2a1, 0x3, 0x2, 0x2, 0x2, 0x2a6, 0x7d, 0x3, 0x2, 0x2, 0x2, 0x2a7, 
       0x2a8, 0x5, 0x80, 0x41, 0x2, 0x2a8, 0x7f, 0x3, 0x2, 0x2, 0x2, 0x2a9, 
       0x2aa, 0x8, 0x41, 0x1, 0x2, 0x2aa, 0x2ab, 0x5, 0x82, 0x42, 0x2, 0x2ab, 
       0x2b1, 0x3, 0x2, 0x2, 0x2, 0x2ac, 0x2ad, 0xc, 0x4, 0x2, 0x2, 0x2ad, 
       0x2ae, 0x7, 0x20, 0x2, 0x2, 0x2ae, 0x2b0, 0x5, 0x82, 0x42, 0x2, 0x2af, 
       0x2ac, 0x3, 0x2, 0x2, 0x2, 0x2b0, 0x2b3, 0x3, 0x2, 0x2, 0x2, 0x2b1, 
       0x2af, 0x3, 0x2, 0x2, 0x2, 0x2b1, 0x2b2, 0x3, 0x2, 0x2, 0x2, 0x2b2, 
       0x81, 0x3, 0x2, 0x2, 0x2, 0x2b3, 0x2b1, 0x3, 0x2, 0x2, 0x2, 0x2b4, 
       0x2b5, 0x8, 0x42, 0x1, 0x2, 0x2b5, 0x2b6, 0x5, 0x84, 0x43, 0x2, 0x2b6, 
       0x2bc, 0x3, 0x2, 0x2, 0x2, 0x2b7, 0x2b8, 0xc, 0x4, 0x2, 0x2, 0x2b8, 
       0x2b9, 0x7, 0x21, 0x2, 0x2, 0x2b9, 0x2bb, 0x5, 0x84, 0x43, 0x2, 0x2ba, 
       0x2b7, 0x3, 0x2, 0x2, 0x2, 0x2bb, 0x2be, 0x3, 0x2, 0x2, 0x2, 0x2bc, 
       0x2ba, 0x3, 0x2, 0x2, 0x2, 0x2bc, 0x2bd, 0x3, 0x2, 0x2, 0x2, 0x2bd, 
       0x83, 0x3, 0x2, 0x2, 0x2, 0x2be, 0x2bc, 0x3, 0x2, 0x2, 0x2, 0x2bf, 
       0x2c0, 0x8, 0x43, 0x1, 0x2, 0x2c0, 0x2c1, 0x5, 0x86, 0x44, 0x2, 0x2c1, 
       0x2ca, 0x3, 0x2, 0x2, 0x2, 0x2c2, 0x2c3, 0xc, 0x4, 0x2, 0x2, 0x2c3, 
       0x2c5, 0x9, 0x2, 0x2, 0x2, 0x2c4, 0x2c6, 0x5, 0x92, 0x4a, 0x2, 0x2c5, 
       0x2c4, 0x3, 0x2, 0x2, 0x2, 0x2c5, 0x2c6, 0x3, 0x2, 0x2, 0x2, 0x2c6, 
       0x2c7, 0x3, 0x2, 0x2, 0x2, 0x2c7, 0x2c9, 0x5, 0x86, 0x44, 0x2, 0x2c8, 
       0x2c2, 0x3, 0x2, 0x2, 0x2, 0x2c9, 0x2cc, 0x3, 0x2, 0x2, 0x2, 0x2ca, 
       0x2c8, 0x3, 0x2, 0x2, 0x2, 0x2ca, 0x2cb, 0x3, 0x2, 0x2, 0x2, 0x2cb, 
       0x85, 0x3, 0x2, 0x2, 0x2, 0x2cc, 0x2ca, 0x3, 0x2, 0x2, 0x2, 0x2cd, 
       0x2ce, 0x8, 0x44, 0x1, 0x2, 0x2ce, 0x2cf, 0x5, 0x88, 0x45, 0x2, 0x2cf, 
       0x2d8, 0x3, 0x2, 0x2, 0x2, 0x2d0, 0x2d1, 0xc, 0x4, 0x2, 0x2, 0x2d1, 
       0x2d3, 0x9, 0x3, 0x2, 0x2, 0x2d2, 0x2d4, 0x5, 0x92, 0x4a, 0x2, 0x2d3, 
       0x2d2, 0x3, 0x2, 0x2, 0x2, 0x2d3, 0x2d4, 0x3, 0x2, 0x2, 0x2, 0x2d4, 
       0x2d5, 0x3, 0x2, 0x2, 0x2, 0x2d5, 0x2d7, 0x5, 0x88, 0x45, 0x2, 0x2d6, 
       0x2d0, 0x3, 0x2, 0x2, 0x2, 0x2d7, 0x2da, 0x3, 0x2, 0x2, 0x2, 0x2d8, 
       0x2d6, 0x3, 0x2, 0x2, 0x2, 0x2d8, 0x2d9, 0x3, 0x2, 0x2, 0x2, 0x2d9, 
       0x87, 0x3, 0x2, 0x2, 0x2, 0x2da, 0x2d8, 0x3, 0x2, 0x2, 0x2, 0x2db, 
       0x2dc, 0x8, 0x45, 0x1, 0x2, 0x2dc, 0x2dd, 0x5, 0x8a, 0x46, 0x2, 0x2dd, 
       0x2e6, 0x3, 0x2, 0x2, 0x2, 0x2de, 0x2df, 0xc, 0x4, 0x2, 0x2, 0x2df, 
       0x2e1, 0x9, 0x4, 0x2, 0x2, 0x2e0, 0x2e2, 0x5, 0x92, 0x4a, 0x2, 0x2e1, 
       0x2e0, 0x3, 0x2, 0x2, 0x2, 0x2e1, 0x2e2, 0x3, 0x2, 0x2, 0x2, 0x2e2, 
       0x2e3, 0x3, 0x2, 0x2, 0x2, 0x2e3, 0x2e5, 0x5, 0x8a, 0x46, 0x2, 0x2e4, 
       0x2de, 0x3, 0x2, 0x2, 0x2, 0x2e5, 0x2e8, 0x3, 0x2, 0x2, 0x2, 0x2e6, 
       0x2e4, 0x3, 0x2, 0x2, 0x2, 0x2e6, 0x2e7, 0x3, 0x2, 0x2, 0x2, 0x2e7, 
       0x89, 0x3, 0x2, 0x2, 0x2, 0x2e8, 0x2e6, 0x3, 0x2, 0x2, 0x2, 0x2e9, 
       0x2ea, 0x8, 0x46, 0x1, 0x2, 0x2ea, 0x2eb, 0x5, 0x8c, 0x47, 0x2, 0x2eb, 
       0x2f4, 0x3, 0x2, 0x2, 0x2, 0x2ec, 0x2ed, 0xc, 0x4, 0x2, 0x2, 0x2ed, 
       0x2ef, 0x9, 0x5, 0x2, 0x2, 0x2ee, 0x2f0, 0x5, 0x92, 0x4a, 0x2, 0x2ef, 
       0x2ee, 0x3, 0x2, 0x2, 0x2, 0x2ef, 0x2f0, 0x3, 0x2, 0x2, 0x2, 0x2f0, 
       0x2f1, 0x3, 0x2, 0x2, 0x2, 0x2f1, 0x2f3, 0x5, 0x8c, 0x47, 0x2, 0x2f2, 
       0x2ec, 0x3, 0x2, 0x2, 0x2, 0x2f3, 0x2f6, 0x3, 0x2, 0x2, 0x2, 0x2f4, 
       0x2f2, 0x3, 0x2, 0x2, 0x2, 0x2f4, 0x2f5, 0x3, 0x2, 0x2, 0x2, 0x2f5, 
       0x8b, 0x3, 0x2, 0x2, 0x2, 0x2f6, 0x2f4, 0x3, 0x2, 0x2, 0x2, 0x2f7, 
       0x2f9, 0x9, 0x6, 0x2, 0x2, 0x2f8, 0x2fa, 0x5, 0x92, 0x4a, 0x2, 0x2f9, 
       0x2f8, 0x3, 0x2, 0x2, 0x2, 0x2f9, 0x2fa, 0x3, 0x2, 0x2, 0x2, 0x2fa, 
       0x2fb, 0x3, 0x2, 0x2, 0x2, 0x2fb, 0x2fe, 0x5, 0x8c, 0x47, 0x2, 0x2fc, 
       0x2fe, 0x5, 0x8e, 0x48, 0x2, 0x2fd, 0x2f7, 0x3, 0x2, 0x2, 0x2, 0x2fd, 
       0x2fc, 0x3, 0x2, 0x2, 0x2, 0x2fe, 0x8d, 0x3, 0x2, 0x2, 0x2, 0x2ff, 
       0x300, 0x8, 0x48, 0x1, 0x2, 0x300, 0x301, 0x5, 0x90, 0x49, 0x2, 0x301, 
       0x306, 0x3, 0x2, 0x2, 0x2, 0x302, 0x303, 0xc, 0x4, 0x2, 0x2, 0x303, 
       0x305, 0x5, 0x94, 0x4b, 0x2, 0x304, 0x302, 0x3, 0x2, 0x2, 0x2, 0x305, 
       0x308, 0x3, 0x2, 0x2, 0x2, 0x306, 0x304, 0x3, 0x2, 0x2, 0x2, 0x306, 
       0x307, 0x3, 0x2, 0x2, 0x2, 0x307, 0x8f, 0x3, 0x2, 0x2, 0x2, 0x308, 
       0x306, 0x3, 0x2, 0x2, 0x2, 0x309, 0x30a, 0x5, 0x96, 0x4c, 0x2, 0x30a, 
       0x30b, 0x5, 0x92, 0x4a, 0x2, 0x30b, 0x30e, 0x3, 0x2, 0x2, 0x2, 0x30c, 
       0x30e, 0x5, 0x96, 0x4c, 0x2, 0x30d, 0x309, 0x3, 0x2, 0x2, 0x2, 0x30d, 
       0x30c, 0x3, 0x2, 0x2, 0x2, 0x30e, 0x91, 0x3, 0x2, 0x2, 0x2, 0x30f, 
       0x310, 0x7, 0x37, 0x2, 0x2, 0x310, 0x311, 0x5, 0x52, 0x2a, 0x2, 0x311, 
       0x312, 0x7, 0x38, 0x2, 0x2, 0x312, 0x318, 0x3, 0x2, 0x2, 0x2, 0x313, 
       0x314, 0x7, 0x37, 0x2, 0x2, 0x314, 0x315, 0x5, 0x3c, 0x1f, 0x2, 0x315, 
       0x316, 0x7, 0x38, 0x2, 0x2, 0x316, 0x318, 0x3, 0x2, 0x2, 0x2, 0x317, 
       0x30f, 0x3, 0x2, 0x2, 0x2, 0x317, 0x313, 0x3, 0x2, 0x2, 0x2, 0x318, 
       0x93, 0x3, 0x2, 0x2, 0x2, 0x319, 0x31a, 0x7, 0x2d, 0x2, 0x2, 0x31a, 
       0x31f, 0x5, 0x98, 0x4d, 0x2, 0x31b, 0x31f, 0x5, 0x98, 0x4d, 0x2, 
       0x31c, 0x31f, 0x7, 0x3b, 0x2, 0x2, 0x31d, 0x31f, 0x7, 0x3c, 0x2, 
       0x2, 0x31e, 0x319, 0x3, 0x2, 0x2, 0x2, 0x31e, 0x31b, 0x3, 0x2, 0x2, 
       0x2, 0x31e, 0x31c, 0x3, 0x2, 0x2, 0x2, 0x31e, 0x31d, 0x3, 0x2, 0x2, 
       0x2, 0x31f, 0x95, 0x3, 0x2, 0x2, 0x2, 0x320, 0x321, 0x7, 0x2a, 0x2, 
       0x2, 0x321, 0x324, 0x5, 0x96, 0x4c, 0x2, 0x322, 0x324, 0x5, 0x94, 
       0x4b, 0x2, 0x323, 0x320, 0x3, 0x2, 0x2, 0x2, 0x323, 0x322, 0x3, 0x2, 
       0x2, 0x2, 0x324, 0x97, 0x3, 0x2, 0x2, 0x2, 0x325, 0x326, 0x8, 0x4d, 
       0x1, 0x2, 0x326, 0x327, 0x5, 0x9a, 0x4e, 0x2, 0x327, 0x330, 0x3, 
       0x2, 0x2, 0x2, 0x328, 0x329, 0xc, 0x5, 0x2, 0x2, 0x329, 0x32a, 0x7, 
       0x2e, 0x2, 0x2, 0x32a, 0x32f, 0x7, 0x3b, 0x2, 0x2, 0x32b, 0x32c, 
       0xc, 0x4, 0x2, 0x2, 0x32c, 0x32d, 0x7, 0x2e, 0x2, 0x2, 0x32d, 0x32f, 
       0x7, 0x3a, 0x2, 0x2, 0x32e, 0x328, 0x3, 0x2, 0x2, 0x2, 0x32e, 0x32b, 
       0x3, 0x2, 0x2, 0x2, 0x32f, 0x332, 0x3, 0x2, 0x2, 0x2, 0x330, 0x32e, 
       0x3, 0x2, 0x2, 0x2, 0x330, 0x331, 0x3, 0x2, 0x2, 0x2, 0x331, 0x99, 
       0x3, 0x2, 0x2, 0x2, 0x332, 0x330, 0x3, 0x2, 0x2, 0x2, 0x333, 0x33e, 
       0x5, 0x9c, 0x4f, 0x2, 0x334, 0x33e, 0x7, 0x3d, 0x2, 0x2, 0x335, 0x33e, 
       0x7, 0x3e, 0x2, 0x2, 0x336, 0x33e, 0x5, 0x9e, 0x50, 0x2, 0x337, 0x33e, 
       0x5, 0xa4, 0x53, 0x2, 0x338, 0x33e, 0x5, 0xa8, 0x55, 0x2, 0x339, 
       0x33a, 0x7, 0x35, 0x2, 0x2, 0x33a, 0x33b, 0x5, 0x5a, 0x2e, 0x2, 0x33b, 
       0x33c, 0x7, 0x36, 0x2, 0x2, 0x33c, 0x33e, 0x3, 0x2, 0x2, 0x2, 0x33d, 
       0x333, 0x3, 0x2, 0x2, 0x2, 0x33d, 0x334, 0x3, 0x2, 0x2, 0x2, 0x33d, 
       0x335, 0x3, 0x2, 0x2, 0x2, 0x33d, 0x336, 0x3, 0x2, 0x2, 0x2, 0x33d, 
       0x337, 0x3, 0x2, 0x2, 0x2, 0x33d, 0x338, 0x3, 0x2, 0x2, 0x2, 0x33d, 
       0x339, 0x3, 0x2, 0x2, 0x2, 0x33e, 0x9b, 0x3, 0x2, 0x2, 0x2, 0x33f, 
       0x340, 0x7, 0x12, 0x2, 0x2, 0x340, 0x341, 0x7, 0x33, 0x2, 0x2, 0x341, 
       0x342, 0x5, 0x3c, 0x1f, 0x2, 0x342, 0x343, 0x5, 0xa2, 0x52, 0x2, 
       0x343, 0x344, 0x7, 0x34, 0x2, 0x2, 0x344, 0x9d, 0x3, 0x2, 0x2, 0x2, 
       0x345, 0x346, 0x7, 0x11, 0x2, 0x2, 0x346, 0x34b, 0x7, 0x33, 0x2, 
       0x2, 0x347, 0x349, 0x5, 0xa0, 0x51, 0x2, 0x348, 0x34a, 0x7, 0x30, 
       0x2, 0x2, 0x349, 0x348, 0x3, 0x2, 0x2, 0x2, 0x349, 0x34a, 0x3, 0x2, 
       0x2, 0x2, 0x34a, 0x34c, 0x3, 0x2, 0x2, 0x2, 0x34b, 0x347, 0x3, 0x2, 
       0x2, 0x2, 0x34b, 0x34c, 0x3, 0x2, 0x2, 0x2, 0x34c, 0x34d, 0x3, 0x2, 
       0x2, 0x2, 0x34d, 0x34e, 0x7, 0x34, 0x2, 0x2, 0x34e, 0x9f, 0x3, 0x2, 
       0x2, 0x2, 0x34f, 0x354, 0x5, 0xa2, 0x52, 0x2, 0x350, 0x351, 0x7, 
       0x30, 0x2, 0x2, 0x351, 0x353, 0x5, 0xa2, 0x52, 0x2, 0x352, 0x350, 
       0x3, 0x2, 0x2, 0x2, 0x353, 0x356, 0x3, 0x2, 0x2, 0x2, 0x354, 0x352, 
       0x3, 0x2, 0x2, 0x2, 0x354, 0x355, 0x3, 0x2, 0x2, 0x2, 0x355, 0xa1, 
       0x3, 0x2, 0x2, 0x2, 0x356, 0x354, 0x3, 0x2, 0x2, 0x2, 0x357, 0x358, 
       0x5, 0xaa, 0x56, 0x2, 0x358, 0x359, 0x7, 0x2f, 0x2, 0x2, 0x359, 0x35a, 
       0x5, 0x5a, 0x2e, 0x2, 0x35a, 0x35f, 0x3, 0x2, 0x2, 0x2, 0x35b, 0x35c, 
       0x7, 0x3b, 0x2, 0x2, 0x35c, 0x35d, 0x7, 0x2f, 0x2, 0x2, 0x35d, 0x35f, 
       0x5, 0x5a, 0x2e, 0x2, 0x35e, 0x357, 0x3, 0x2, 0x2, 0x2, 0x35e, 0x35b, 
       0x3, 0x2, 0x2, 0x2, 0x35f, 0xa3, 0x3, 0x2, 0x2, 0x2, 0x360, 0x361, 
       0x7, 0x35, 0x2, 0x2, 0x361, 0x367, 0x7, 0x36, 0x2, 0x2, 0x362, 0x363, 
       0x7, 0x35, 0x2, 0x2, 0x363, 0x364, 0x5, 0xa6, 0x54, 0x2, 0x364, 0x365, 
       0x7, 0x36, 0x2, 0x2, 0x365, 0x367, 0x3, 0x2, 0x2, 0x2, 0x366, 0x360, 
       0x3, 0x2, 0x2, 0x2, 0x366, 0x362, 0x3, 0x2, 0x2, 0x2, 0x367, 0xa5, 
       0x3, 0x2, 0x2, 0x2, 0x368, 0x36b, 0x5, 0x5a, 0x2e, 0x2, 0x369, 0x36a, 
       0x7, 0x30, 0x2, 0x2, 0x36a, 0x36c, 0x5, 0x5a, 0x2e, 0x2, 0x36b, 0x369, 
       0x3, 0x2, 0x2, 0x2, 0x36c, 0x36d, 0x3, 0x2, 0x2, 0x2, 0x36d, 0x36b, 
       0x3, 0x2, 0x2, 0x2, 0x36d, 0x36e, 0x3, 0x2, 0x2, 0x2, 0x36e, 0xa7, 
       0x3, 0x2, 0x2, 0x2, 0x36f, 0x370, 0x5, 0xaa, 0x56, 0x2, 0x370, 0xa9, 
       0x3, 0x2, 0x2, 0x2, 0x371, 0x376, 0x5, 0xac, 0x57, 0x2, 0x372, 0x373, 
       0x7, 0x35, 0x2, 0x2, 0x373, 0x374, 0x9, 0x7, 0x2, 0x2, 0x374, 0x376, 
       0x7, 0x36, 0x2, 0x2, 0x375, 0x371, 0x3, 0x2, 0x2, 0x2, 0x375, 0x372, 
       0x3, 0x2, 0x2, 0x2, 0x376, 0xab, 0x3, 0x2, 0x2, 0x2, 0x377, 0x37c, 
       0x7, 0x3a, 0x2, 0x2, 0x378, 0x379, 0x7, 0x1c, 0x2, 0x2, 0x379, 0x37b, 
       0x9, 0x8, 0x2, 0x2, 0x37a, 0x378, 0x3, 0x2, 0x2, 0x2, 0x37b, 0x37e, 
       0x3, 0x2, 0x2, 0x2, 0x37c, 0x37a, 0x3, 0x2, 0x2, 0x2, 0x37c, 0x37d, 
       0x3, 0x2, 0x2, 0x2, 0x37d, 0xad, 0x3, 0x2, 0x2, 0x2, 0x37e, 0x37c, 
       0x3, 0x2, 0x2, 0x2, 0x62, 0xb0, 0xb3, 0xb6, 0xb9, 0xbf, 0xc2, 0xc5, 
       0xc8, 0xcc, 0xde, 0xec, 0xf3, 0xfc, 0x105, 0x10e, 0x116, 0x120, 0x125, 
       0x12b, 0x138, 0x141, 0x148, 0x14f, 0x154, 0x15b, 0x160, 0x16d, 0x170, 
       0x17b, 0x17f, 0x184, 0x18d, 0x195, 0x197, 0x19f, 0x1a7, 0x1b1, 0x1bb, 
       0x1c0, 0x1c5, 0x1ce, 0x1dc, 0x1ef, 0x1f5, 0x1f7, 0x200, 0x20d, 0x214, 
       0x21a, 0x21c, 0x225, 0x22d, 0x232, 0x23a, 0x23e, 0x24c, 0x24e, 0x256, 
       0x260, 0x269, 0x276, 0x280, 0x283, 0x288, 0x28c, 0x29c, 0x2a3, 0x2a5, 
       0x2b1, 0x2bc, 0x2c5, 0x2ca, 0x2d3, 0x2d8, 0x2e1, 0x2e6, 0x2ef, 0x2f4, 
       0x2f9, 0x2fd, 0x306, 0x30d, 0x317, 0x31e, 0x323, 0x32e, 0x330, 0x33d, 
       0x349, 0x34b, 0x354, 0x35e, 0x366, 0x36d, 0x375, 0x37c, 
  };

  _serializedATN.insert(_serializedATN.end(), serializedATNSegment0,
    serializedATNSegment0 + sizeof(serializedATNSegment0) / sizeof(serializedATNSegment0[0]));


  atn::ATNDeserializer deserializer;
  _atn = deserializer.deserialize(_serializedATN);

  size_t count = _atn.getNumberOfDecisions();
  _decisionToDFA.reserve(count);
  for (size_t i = 0; i < count; i++) { 
    _decisionToDFA.emplace_back(_atn.getDecisionState(i), i);
  }
}

bapelParser::Initializer bapelParser::_init;
