
// Generated from cpp_parser/bapel.g4 by ANTLR 4.9.2

#pragma once


#include "antlr4-runtime.h"




class  bapelParser : public antlr4::Parser {
public:
  enum {
    WORKSPACE = 1, PACKAGES = 2, PREFIX = 3, MODULE = 4, IMPLEMENTS = 5, 
    IMPORTS = 6, IMPLS = 7, FLAGS = 8, IN = 9, PUB = 10, DECL = 11, FN = 12, 
    TYPE = 13, FORALL = 14, STRUCT = 15, VARIANT = 16, MATCH = 17, SET = 18, 
    LET = 19, RETURN = 20, IF = 21, ELSE = 22, FOR = 23, TRAIT = 24, IMPL = 25, 
    DOUBLE_COLON = 26, ARROW = 27, FAT_ARROW = 28, LARROW = 29, OR = 30, 
    AND = 31, NE = 32, EQ = 33, GE = 34, LE = 35, GT = 36, LT = 37, PLUS = 38, 
    MINUS = 39, MUL = 40, DIV = 41, NOT = 42, AMP = 43, DOT = 44, ASSIGN = 45, 
    COMMA = 46, SEMI = 47, COLON = 48, LBRACE = 49, RBRACE = 50, LPAREN = 51, 
    RPAREN = 52, LBRACKET = 53, RBRACKET = 54, SINGLE_QUOTE = 55, IDENTIFIER = 56, 
    INT_LITERAL = 57, FLOAT_LITERAL = 58, RUNE_LITERAL = 59, STRING_LITERAL = 60, 
    RAW_STRING_LITERAL = 61, UNTERMINATED_RUNE_LITERAL = 62, UNTERMINATED_STRING_LITERAL = 63, 
    UNTERMINATED_RAW_STRING_LITERAL = 64, UNTERMINATED_BLOCK_COMMENT = 65, 
    WS = 66, LINE_COMMENT = 67, BLOCK_COMMENT = 68
  };

  enum {
    RuleSourceFile = 0, RuleModuleHeader = 1, RuleImplementsHeader = 2, 
    RuleWorkspace = 3, RulePackagesSection = 4, RulePackageRule = 5, RuleImportsSection = 6, 
    RuleImplsSection = 7, RuleFlagsSection = 8, RuleModuleID = 9, RuleFilename = 10, 
    RuleSources = 11, RuleSource = 12, RuleTraitDecl = 13, RuleTraitMethod = 14, 
    RuleImplBlock = 15, RuleDeclNoExport = 16, RuleFunctionNoExport = 17, 
    RuleFunctionArgs = 18, RuleArg = 19, RuleDecl = 20, RuleUnexportedDecl = 21, 
    RuleDeclNoTerm = 22, RuleTermDecl = 23, RuleTypeDecl = 24, RuleTypeAbstraction = 25, 
    RuleBoundedTvar = 26, RuleTvar = 27, RuleTraitBound = 28, RuleType_ = 29, 
    RuleForallType = 30, RuleFunctionType = 31, RulePtrType = 32, RuleAppType = 33, 
    RulePrimaryType = 34, RuleArrayType = 35, RuleStructType = 36, RuleFields = 37, 
    RuleField = 38, RuleTupleType = 39, RuleTupleTypeArgs = 40, RuleVariantType = 41, 
    RuleTags = 42, RuleTag = 43, RuleExpression = 44, RuleExpressionWithoutBlock = 45, 
    RuleExpressionWithBlock = 46, RuleAssignTerm = 47, RuleReturnTerm = 48, 
    RuleIfTerm = 49, RuleForTerm = 50, RuleLambdaTerm = 51, RuleMatchTerm = 52, 
    RuleMatchArms = 53, RuleMatchArm = 54, RuleSetTerm = 55, RuleBlockExpr = 56, 
    RuleBlockStatements = 57, RuleStatements = 58, RuleStatement = 59, RuleLetStatement = 60, 
    RuleExpressionStatement = 61, RuleOperatorExpr = 62, RuleLogicalOrExpr = 63, 
    RuleLogicalAndExpr = 64, RuleEqualityExpr = 65, RuleComparisonExpr = 66, 
    RuleAdditiveExpr = 67, RuleMultiplicativeExpr = 68, RuleUnaryExpr = 69, 
    RuleApplicativeExpr = 70, RuleTypeApplicativeExpr = 71, RuleTypeApplicativeArgs = 72, 
    RuleBasePrimaryExpr = 73, RulePrimaryExpr = 74, RuleProjectionExpr = 75, 
    RuleDerefExpr = 76, RuleInjectionExpr = 77, RuleStructExpr = 78, RuleLabelValues = 79, 
    RuleLabelValue = 80, RuleTupleExpr = 81, RuleTupleExprArgs = 82, RuleVarExpr = 83, 
    RuleId = 84, RuleIdTokens = 85
  };

  explicit bapelParser(antlr4::TokenStream *input);
  ~bapelParser();

  virtual std::string getGrammarFileName() const override;
  virtual const antlr4::atn::ATN& getATN() const override { return _atn; };
  virtual const std::vector<std::string>& getTokenNames() const override { return _tokenNames; }; // deprecated: use vocabulary instead.
  virtual const std::vector<std::string>& getRuleNames() const override;
  virtual antlr4::dfa::Vocabulary& getVocabulary() const override;


  class SourceFileContext;
  class ModuleHeaderContext;
  class ImplementsHeaderContext;
  class WorkspaceContext;
  class PackagesSectionContext;
  class PackageRuleContext;
  class ImportsSectionContext;
  class ImplsSectionContext;
  class FlagsSectionContext;
  class ModuleIDContext;
  class FilenameContext;
  class SourcesContext;
  class SourceContext;
  class TraitDeclContext;
  class TraitMethodContext;
  class ImplBlockContext;
  class DeclNoExportContext;
  class FunctionNoExportContext;
  class FunctionArgsContext;
  class ArgContext;
  class DeclContext;
  class UnexportedDeclContext;
  class DeclNoTermContext;
  class TermDeclContext;
  class TypeDeclContext;
  class TypeAbstractionContext;
  class BoundedTvarContext;
  class TvarContext;
  class TraitBoundContext;
  class Type_Context;
  class ForallTypeContext;
  class FunctionTypeContext;
  class PtrTypeContext;
  class AppTypeContext;
  class PrimaryTypeContext;
  class ArrayTypeContext;
  class StructTypeContext;
  class FieldsContext;
  class FieldContext;
  class TupleTypeContext;
  class TupleTypeArgsContext;
  class VariantTypeContext;
  class TagsContext;
  class TagContext;
  class ExpressionContext;
  class ExpressionWithoutBlockContext;
  class ExpressionWithBlockContext;
  class AssignTermContext;
  class ReturnTermContext;
  class IfTermContext;
  class ForTermContext;
  class LambdaTermContext;
  class MatchTermContext;
  class MatchArmsContext;
  class MatchArmContext;
  class SetTermContext;
  class BlockExprContext;
  class BlockStatementsContext;
  class StatementsContext;
  class StatementContext;
  class LetStatementContext;
  class ExpressionStatementContext;
  class OperatorExprContext;
  class LogicalOrExprContext;
  class LogicalAndExprContext;
  class EqualityExprContext;
  class ComparisonExprContext;
  class AdditiveExprContext;
  class MultiplicativeExprContext;
  class UnaryExprContext;
  class ApplicativeExprContext;
  class TypeApplicativeExprContext;
  class TypeApplicativeArgsContext;
  class BasePrimaryExprContext;
  class PrimaryExprContext;
  class ProjectionExprContext;
  class DerefExprContext;
  class InjectionExprContext;
  class StructExprContext;
  class LabelValuesContext;
  class LabelValueContext;
  class TupleExprContext;
  class TupleExprArgsContext;
  class VarExprContext;
  class IdContext;
  class IdTokensContext; 

  class  SourceFileContext : public antlr4::ParserRuleContext {
  public:
    SourceFileContext(antlr4::ParserRuleContext *parent, size_t invokingState);
   
    SourceFileContext() = default;
    void copyFrom(SourceFileContext *context);
    using antlr4::ParserRuleContext::copyFrom;

    virtual size_t getRuleIndex() const override;

   
  };

  class  ImplSourceFileContext : public SourceFileContext {
  public:
    ImplSourceFileContext(SourceFileContext *ctx);

    ImplementsHeaderContext *implementsHeader();
    antlr4::tree::TerminalNode *EOF();
    ImportsSectionContext *importsSection();
    ImplsSectionContext *implsSection();
    FlagsSectionContext *flagsSection();
    SourcesContext *sources();

    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
  };

  class  BaseSourceFileContext : public SourceFileContext {
  public:
    BaseSourceFileContext(SourceFileContext *ctx);

    ModuleHeaderContext *moduleHeader();
    antlr4::tree::TerminalNode *EOF();
    ImportsSectionContext *importsSection();
    ImplsSectionContext *implsSection();
    FlagsSectionContext *flagsSection();
    SourcesContext *sources();

    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
  };

  SourceFileContext* sourceFile();

  class  ModuleHeaderContext : public antlr4::ParserRuleContext {
  public:
    ModuleHeaderContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *MODULE();
    ModuleIDContext *moduleID();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ModuleHeaderContext* moduleHeader();

  class  ImplementsHeaderContext : public antlr4::ParserRuleContext {
  public:
    ImplementsHeaderContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *IMPLEMENTS();
    ModuleIDContext *moduleID();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ImplementsHeaderContext* implementsHeader();

  class  WorkspaceContext : public antlr4::ParserRuleContext {
  public:
    WorkspaceContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *WORKSPACE();
    antlr4::tree::TerminalNode *LBRACE();
    PackagesSectionContext *packagesSection();
    antlr4::tree::TerminalNode *RBRACE();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  WorkspaceContext* workspace();

  class  PackagesSectionContext : public antlr4::ParserRuleContext {
  public:
    PackagesSectionContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *PACKAGES();
    antlr4::tree::TerminalNode *LBRACE();
    antlr4::tree::TerminalNode *RBRACE();
    std::vector<PackageRuleContext *> packageRule();
    PackageRuleContext* packageRule(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  PackagesSectionContext* packagesSection();

  class  PackageRuleContext : public antlr4::ParserRuleContext {
  public:
    PackageRuleContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *PREFIX();
    ModuleIDContext *moduleID();
    antlr4::tree::TerminalNode *IN();
    FilenameContext *filename();
    antlr4::tree::TerminalNode *MODULE();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  PackageRuleContext* packageRule();

  class  ImportsSectionContext : public antlr4::ParserRuleContext {
  public:
    ImportsSectionContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *IMPORTS();
    antlr4::tree::TerminalNode *LBRACE();
    antlr4::tree::TerminalNode *RBRACE();
    std::vector<ModuleIDContext *> moduleID();
    ModuleIDContext* moduleID(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ImportsSectionContext* importsSection();

  class  ImplsSectionContext : public antlr4::ParserRuleContext {
  public:
    ImplsSectionContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *IMPLS();
    antlr4::tree::TerminalNode *LBRACE();
    antlr4::tree::TerminalNode *RBRACE();
    std::vector<FilenameContext *> filename();
    FilenameContext* filename(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ImplsSectionContext* implsSection();

  class  FlagsSectionContext : public antlr4::ParserRuleContext {
  public:
    FlagsSectionContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *FLAGS();
    antlr4::tree::TerminalNode *LBRACE();
    antlr4::tree::TerminalNode *RBRACE();
    std::vector<FilenameContext *> filename();
    FilenameContext* filename(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  FlagsSectionContext* flagsSection();

  class  ModuleIDContext : public antlr4::ParserRuleContext {
  public:
    ModuleIDContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    std::vector<antlr4::tree::TerminalNode *> IDENTIFIER();
    antlr4::tree::TerminalNode* IDENTIFIER(size_t i);
    std::vector<antlr4::tree::TerminalNode *> DOT();
    antlr4::tree::TerminalNode* DOT(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ModuleIDContext* moduleID();

  class  FilenameContext : public antlr4::ParserRuleContext {
  public:
    FilenameContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *STRING_LITERAL();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  FilenameContext* filename();

  class  SourcesContext : public antlr4::ParserRuleContext {
  public:
    SourcesContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    std::vector<SourceContext *> source();
    SourceContext* source(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  SourcesContext* sources();

  class  SourceContext : public antlr4::ParserRuleContext {
  public:
    SourceContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    DeclNoExportContext *declNoExport();
    FunctionNoExportContext *functionNoExport();
    antlr4::tree::TerminalNode *PUB();
    TraitDeclContext *traitDecl();
    ImplBlockContext *implBlock();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  SourceContext* source();

  class  TraitDeclContext : public antlr4::ParserRuleContext {
  public:
    TraitDeclContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *TRAIT();
    IdContext *id();
    antlr4::tree::TerminalNode *LBRACE();
    antlr4::tree::TerminalNode *RBRACE();
    TypeAbstractionContext *typeAbstraction();
    std::vector<TraitMethodContext *> traitMethod();
    TraitMethodContext* traitMethod(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TraitDeclContext* traitDecl();

  class  TraitMethodContext : public antlr4::ParserRuleContext {
  public:
    TraitMethodContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *FN();
    IdContext *id();
    FunctionArgsContext *functionArgs();
    antlr4::tree::TerminalNode *ARROW();
    Type_Context *type_();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TraitMethodContext* traitMethod();

  class  ImplBlockContext : public antlr4::ParserRuleContext {
  public:
    ImplBlockContext(antlr4::ParserRuleContext *parent, size_t invokingState);
   
    ImplBlockContext() = default;
    void copyFrom(ImplBlockContext *context);
    using antlr4::ParserRuleContext::copyFrom;

    virtual size_t getRuleIndex() const override;

   
  };

  class  InherentImplContext : public ImplBlockContext {
  public:
    InherentImplContext(ImplBlockContext *ctx);

    antlr4::tree::TerminalNode *IMPL();
    Type_Context *type_();
    antlr4::tree::TerminalNode *LBRACE();
    antlr4::tree::TerminalNode *RBRACE();
    TypeAbstractionContext *typeAbstraction();
    std::vector<FunctionNoExportContext *> functionNoExport();
    FunctionNoExportContext* functionNoExport(size_t i);

    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
  };

  class  TraitImplContext : public ImplBlockContext {
  public:
    TraitImplContext(ImplBlockContext *ctx);

    antlr4::tree::TerminalNode *IMPL();
    std::vector<Type_Context *> type_();
    Type_Context* type_(size_t i);
    antlr4::tree::TerminalNode *FOR();
    antlr4::tree::TerminalNode *LBRACE();
    antlr4::tree::TerminalNode *RBRACE();
    TypeAbstractionContext *typeAbstraction();
    std::vector<FunctionNoExportContext *> functionNoExport();
    FunctionNoExportContext* functionNoExport(size_t i);

    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
  };

  ImplBlockContext* implBlock();

  class  DeclNoExportContext : public antlr4::ParserRuleContext {
  public:
    DeclNoExportContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    DeclNoTermContext *declNoTerm();
    antlr4::tree::TerminalNode *DECL();
    TermDeclContext *termDecl();
    antlr4::tree::TerminalNode *PUB();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  DeclNoExportContext* declNoExport();

  class  FunctionNoExportContext : public antlr4::ParserRuleContext {
  public:
    FunctionNoExportContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *FN();
    IdContext *id();
    FunctionArgsContext *functionArgs();
    antlr4::tree::TerminalNode *ARROW();
    Type_Context *type_();
    BlockExprContext *blockExpr();
    TypeAbstractionContext *typeAbstraction();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  FunctionNoExportContext* functionNoExport();

  class  FunctionArgsContext : public antlr4::ParserRuleContext {
  public:
    FunctionArgsContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *LPAREN();
    antlr4::tree::TerminalNode *RPAREN();
    std::vector<ArgContext *> arg();
    ArgContext* arg(size_t i);
    std::vector<antlr4::tree::TerminalNode *> COMMA();
    antlr4::tree::TerminalNode* COMMA(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  FunctionArgsContext* functionArgs();

  class  ArgContext : public antlr4::ParserRuleContext {
  public:
    ArgContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *IDENTIFIER();
    antlr4::tree::TerminalNode *COLON();
    Type_Context *type_();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ArgContext* arg();

  class  DeclContext : public antlr4::ParserRuleContext {
  public:
    DeclContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *PUB();
    UnexportedDeclContext *unexportedDecl();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  DeclContext* decl();

  class  UnexportedDeclContext : public antlr4::ParserRuleContext {
  public:
    UnexportedDeclContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    TermDeclContext *termDecl();
    TypeDeclContext *typeDecl();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  UnexportedDeclContext* unexportedDecl();

  class  DeclNoTermContext : public antlr4::ParserRuleContext {
  public:
    DeclNoTermContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *PUB();
    TypeDeclContext *typeDecl();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  DeclNoTermContext* declNoTerm();

  class  TermDeclContext : public antlr4::ParserRuleContext {
  public:
    TermDeclContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    IdContext *id();
    antlr4::tree::TerminalNode *COLON();
    Type_Context *type_();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TermDeclContext* termDecl();

  class  TypeDeclContext : public antlr4::ParserRuleContext {
  public:
    TypeDeclContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *TYPE();
    IdContext *id();
    antlr4::tree::TerminalNode *ASSIGN();
    Type_Context *type_();
    TypeAbstractionContext *typeAbstraction();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TypeDeclContext* typeDecl();

  class  TypeAbstractionContext : public antlr4::ParserRuleContext {
  public:
    TypeAbstractionContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *LBRACKET();
    std::vector<BoundedTvarContext *> boundedTvar();
    BoundedTvarContext* boundedTvar(size_t i);
    antlr4::tree::TerminalNode *RBRACKET();
    std::vector<antlr4::tree::TerminalNode *> COMMA();
    antlr4::tree::TerminalNode* COMMA(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TypeAbstractionContext* typeAbstraction();

  class  BoundedTvarContext : public antlr4::ParserRuleContext {
  public:
    BoundedTvarContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    TvarContext *tvar();
    antlr4::tree::TerminalNode *COLON();
    TraitBoundContext *traitBound();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  BoundedTvarContext* boundedTvar();

  class  TvarContext : public antlr4::ParserRuleContext {
  public:
    TvarContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *SINGLE_QUOTE();
    antlr4::tree::TerminalNode *IDENTIFIER();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TvarContext* tvar();

  class  TraitBoundContext : public antlr4::ParserRuleContext {
  public:
    TraitBoundContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    std::vector<Type_Context *> type_();
    Type_Context* type_(size_t i);
    std::vector<antlr4::tree::TerminalNode *> PLUS();
    antlr4::tree::TerminalNode* PLUS(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TraitBoundContext* traitBound();

  class  Type_Context : public antlr4::ParserRuleContext {
  public:
    Type_Context(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    ForallTypeContext *forallType();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  Type_Context* type_();

  class  ForallTypeContext : public antlr4::ParserRuleContext {
  public:
    ForallTypeContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *FORALL();
    TypeAbstractionContext *typeAbstraction();
    FunctionTypeContext *functionType();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ForallTypeContext* forallType();

  class  FunctionTypeContext : public antlr4::ParserRuleContext {
  public:
    FunctionTypeContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    PtrTypeContext *ptrType();
    antlr4::tree::TerminalNode *ARROW();
    FunctionTypeContext *functionType();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  FunctionTypeContext* functionType();

  class  PtrTypeContext : public antlr4::ParserRuleContext {
  public:
    PtrTypeContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *AMP();
    PtrTypeContext *ptrType();
    AppTypeContext *appType();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  PtrTypeContext* ptrType();

  class  AppTypeContext : public antlr4::ParserRuleContext {
  public:
    AppTypeContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    PrimaryTypeContext *primaryType();
    AppTypeContext *appType();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  AppTypeContext* appType();
  AppTypeContext* appType(int precedence);
  class  PrimaryTypeContext : public antlr4::ParserRuleContext {
  public:
    PrimaryTypeContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    ArrayTypeContext *arrayType();
    StructTypeContext *structType();
    TupleTypeContext *tupleType();
    VariantTypeContext *variantType();
    antlr4::tree::TerminalNode *SINGLE_QUOTE();
    antlr4::tree::TerminalNode *IDENTIFIER();
    IdContext *id();
    antlr4::tree::TerminalNode *LPAREN();
    Type_Context *type_();
    antlr4::tree::TerminalNode *RPAREN();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  PrimaryTypeContext* primaryType();

  class  ArrayTypeContext : public antlr4::ParserRuleContext {
  public:
    ArrayTypeContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *LBRACKET();
    Type_Context *type_();
    antlr4::tree::TerminalNode *COMMA();
    antlr4::tree::TerminalNode *INT_LITERAL();
    antlr4::tree::TerminalNode *RBRACKET();
    antlr4::tree::TerminalNode *MINUS();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ArrayTypeContext* arrayType();

  class  StructTypeContext : public antlr4::ParserRuleContext {
  public:
    StructTypeContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *STRUCT();
    antlr4::tree::TerminalNode *LBRACE();
    antlr4::tree::TerminalNode *RBRACE();
    FieldsContext *fields();
    antlr4::tree::TerminalNode *COMMA();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  StructTypeContext* structType();

  class  FieldsContext : public antlr4::ParserRuleContext {
  public:
    FieldsContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    std::vector<FieldContext *> field();
    FieldContext* field(size_t i);
    std::vector<antlr4::tree::TerminalNode *> COMMA();
    antlr4::tree::TerminalNode* COMMA(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  FieldsContext* fields();

  class  FieldContext : public antlr4::ParserRuleContext {
  public:
    FieldContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    IdContext *id();
    antlr4::tree::TerminalNode *COLON();
    Type_Context *type_();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  FieldContext* field();

  class  TupleTypeContext : public antlr4::ParserRuleContext {
  public:
    TupleTypeContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *LPAREN();
    antlr4::tree::TerminalNode *RPAREN();
    TupleTypeArgsContext *tupleTypeArgs();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TupleTypeContext* tupleType();

  class  TupleTypeArgsContext : public antlr4::ParserRuleContext {
  public:
    TupleTypeArgsContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    std::vector<Type_Context *> type_();
    Type_Context* type_(size_t i);
    std::vector<antlr4::tree::TerminalNode *> COMMA();
    antlr4::tree::TerminalNode* COMMA(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TupleTypeArgsContext* tupleTypeArgs();

  class  VariantTypeContext : public antlr4::ParserRuleContext {
  public:
    VariantTypeContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *VARIANT();
    antlr4::tree::TerminalNode *LBRACE();
    antlr4::tree::TerminalNode *RBRACE();
    TagsContext *tags();
    antlr4::tree::TerminalNode *COMMA();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  VariantTypeContext* variantType();

  class  TagsContext : public antlr4::ParserRuleContext {
  public:
    TagsContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    std::vector<TagContext *> tag();
    TagContext* tag(size_t i);
    std::vector<antlr4::tree::TerminalNode *> COMMA();
    antlr4::tree::TerminalNode* COMMA(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TagsContext* tags();

  class  TagContext : public antlr4::ParserRuleContext {
  public:
    TagContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    IdContext *id();
    Type_Context *type_();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TagContext* tag();

  class  ExpressionContext : public antlr4::ParserRuleContext {
  public:
    ExpressionContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    ExpressionWithoutBlockContext *expressionWithoutBlock();
    ExpressionWithBlockContext *expressionWithBlock();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ExpressionContext* expression();

  class  ExpressionWithoutBlockContext : public antlr4::ParserRuleContext {
  public:
    ExpressionWithoutBlockContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    AssignTermContext *assignTerm();
    OperatorExprContext *operatorExpr();
    ReturnTermContext *returnTerm();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ExpressionWithoutBlockContext* expressionWithoutBlock();

  class  ExpressionWithBlockContext : public antlr4::ParserRuleContext {
  public:
    ExpressionWithBlockContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    BlockExprContext *blockExpr();
    IfTermContext *ifTerm();
    ForTermContext *forTerm();
    LambdaTermContext *lambdaTerm();
    MatchTermContext *matchTerm();
    SetTermContext *setTerm();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ExpressionWithBlockContext* expressionWithBlock();

  class  AssignTermContext : public antlr4::ParserRuleContext {
  public:
    AssignTermContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *LARROW();
    ExpressionContext *expression();
    IdContext *id();
    TupleExprContext *tupleExpr();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  AssignTermContext* assignTerm();

  class  ReturnTermContext : public antlr4::ParserRuleContext {
  public:
    ReturnTermContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *RETURN();
    ExpressionWithoutBlockContext *expressionWithoutBlock();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ReturnTermContext* returnTerm();

  class  IfTermContext : public antlr4::ParserRuleContext {
  public:
    IfTermContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *IF();
    ExpressionWithoutBlockContext *expressionWithoutBlock();
    std::vector<BlockExprContext *> blockExpr();
    BlockExprContext* blockExpr(size_t i);
    antlr4::tree::TerminalNode *ELSE();
    IfTermContext *ifTerm();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  IfTermContext* ifTerm();

  class  ForTermContext : public antlr4::ParserRuleContext {
  public:
    ForTermContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *FOR();
    ExpressionWithoutBlockContext *expressionWithoutBlock();
    BlockExprContext *blockExpr();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ForTermContext* forTerm();

  class  LambdaTermContext : public antlr4::ParserRuleContext {
  public:
    LambdaTermContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *FN();
    FunctionArgsContext *functionArgs();
    BlockExprContext *blockExpr();
    TypeAbstractionContext *typeAbstraction();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  LambdaTermContext* lambdaTerm();

  class  MatchTermContext : public antlr4::ParserRuleContext {
  public:
    MatchTermContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *MATCH();
    ExpressionContext *expression();
    antlr4::tree::TerminalNode *LBRACE();
    MatchArmsContext *matchArms();
    antlr4::tree::TerminalNode *RBRACE();
    antlr4::tree::TerminalNode *COMMA();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  MatchTermContext* matchTerm();

  class  MatchArmsContext : public antlr4::ParserRuleContext {
  public:
    MatchArmsContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    std::vector<MatchArmContext *> matchArm();
    MatchArmContext* matchArm(size_t i);
    std::vector<antlr4::tree::TerminalNode *> COMMA();
    antlr4::tree::TerminalNode* COMMA(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  MatchArmsContext* matchArms();

  class  MatchArmContext : public antlr4::ParserRuleContext {
  public:
    MatchArmContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    IdContext *id();
    antlr4::tree::TerminalNode *IDENTIFIER();
    antlr4::tree::TerminalNode *FAT_ARROW();
    ExpressionContext *expression();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  MatchArmContext* matchArm();

  class  SetTermContext : public antlr4::ParserRuleContext {
  public:
    SetTermContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *SET();
    ExpressionContext *expression();
    antlr4::tree::TerminalNode *LBRACE();
    LabelValuesContext *labelValues();
    antlr4::tree::TerminalNode *RBRACE();
    antlr4::tree::TerminalNode *COMMA();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  SetTermContext* setTerm();

  class  BlockExprContext : public antlr4::ParserRuleContext {
  public:
    BlockExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *LBRACE();
    BlockStatementsContext *blockStatements();
    antlr4::tree::TerminalNode *RBRACE();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  BlockExprContext* blockExpr();

  class  BlockStatementsContext : public antlr4::ParserRuleContext {
  public:
    BlockStatementsContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    StatementsContext *statements();
    ExpressionWithoutBlockContext *expressionWithoutBlock();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  BlockStatementsContext* blockStatements();

  class  StatementsContext : public antlr4::ParserRuleContext {
  public:
    StatementsContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    std::vector<StatementContext *> statement();
    StatementContext* statement(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  StatementsContext* statements();

  class  StatementContext : public antlr4::ParserRuleContext {
  public:
    StatementContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    LetStatementContext *letStatement();
    ExpressionStatementContext *expressionStatement();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  StatementContext* statement();

  class  LetStatementContext : public antlr4::ParserRuleContext {
  public:
    LetStatementContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *LET();
    IdContext *id();
    antlr4::tree::TerminalNode *COLON();
    Type_Context *type_();
    antlr4::tree::TerminalNode *ASSIGN();
    ExpressionContext *expression();
    antlr4::tree::TerminalNode *SEMI();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  LetStatementContext* letStatement();

  class  ExpressionStatementContext : public antlr4::ParserRuleContext {
  public:
    ExpressionStatementContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    ExpressionWithoutBlockContext *expressionWithoutBlock();
    antlr4::tree::TerminalNode *SEMI();
    ExpressionWithBlockContext *expressionWithBlock();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ExpressionStatementContext* expressionStatement();

  class  OperatorExprContext : public antlr4::ParserRuleContext {
  public:
    OperatorExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    LogicalOrExprContext *logicalOrExpr();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  OperatorExprContext* operatorExpr();

  class  LogicalOrExprContext : public antlr4::ParserRuleContext {
  public:
    LogicalOrExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    LogicalAndExprContext *logicalAndExpr();
    LogicalOrExprContext *logicalOrExpr();
    antlr4::tree::TerminalNode *OR();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  LogicalOrExprContext* logicalOrExpr();
  LogicalOrExprContext* logicalOrExpr(int precedence);
  class  LogicalAndExprContext : public antlr4::ParserRuleContext {
  public:
    LogicalAndExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    EqualityExprContext *equalityExpr();
    LogicalAndExprContext *logicalAndExpr();
    antlr4::tree::TerminalNode *AND();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  LogicalAndExprContext* logicalAndExpr();
  LogicalAndExprContext* logicalAndExpr(int precedence);
  class  EqualityExprContext : public antlr4::ParserRuleContext {
  public:
    EqualityExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    ComparisonExprContext *comparisonExpr();
    EqualityExprContext *equalityExpr();
    antlr4::tree::TerminalNode *NE();
    antlr4::tree::TerminalNode *EQ();
    TypeApplicativeArgsContext *typeApplicativeArgs();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  EqualityExprContext* equalityExpr();
  EqualityExprContext* equalityExpr(int precedence);
  class  ComparisonExprContext : public antlr4::ParserRuleContext {
  public:
    ComparisonExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    AdditiveExprContext *additiveExpr();
    ComparisonExprContext *comparisonExpr();
    antlr4::tree::TerminalNode *GT();
    antlr4::tree::TerminalNode *GE();
    antlr4::tree::TerminalNode *LT();
    antlr4::tree::TerminalNode *LE();
    TypeApplicativeArgsContext *typeApplicativeArgs();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ComparisonExprContext* comparisonExpr();
  ComparisonExprContext* comparisonExpr(int precedence);
  class  AdditiveExprContext : public antlr4::ParserRuleContext {
  public:
    AdditiveExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    MultiplicativeExprContext *multiplicativeExpr();
    AdditiveExprContext *additiveExpr();
    antlr4::tree::TerminalNode *PLUS();
    antlr4::tree::TerminalNode *MINUS();
    TypeApplicativeArgsContext *typeApplicativeArgs();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  AdditiveExprContext* additiveExpr();
  AdditiveExprContext* additiveExpr(int precedence);
  class  MultiplicativeExprContext : public antlr4::ParserRuleContext {
  public:
    MultiplicativeExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    UnaryExprContext *unaryExpr();
    MultiplicativeExprContext *multiplicativeExpr();
    antlr4::tree::TerminalNode *MUL();
    antlr4::tree::TerminalNode *DIV();
    TypeApplicativeArgsContext *typeApplicativeArgs();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  MultiplicativeExprContext* multiplicativeExpr();
  MultiplicativeExprContext* multiplicativeExpr(int precedence);
  class  UnaryExprContext : public antlr4::ParserRuleContext {
  public:
    UnaryExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    UnaryExprContext *unaryExpr();
    antlr4::tree::TerminalNode *NOT();
    antlr4::tree::TerminalNode *MINUS();
    TypeApplicativeArgsContext *typeApplicativeArgs();
    ApplicativeExprContext *applicativeExpr();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  UnaryExprContext* unaryExpr();

  class  ApplicativeExprContext : public antlr4::ParserRuleContext {
  public:
    ApplicativeExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    TypeApplicativeExprContext *typeApplicativeExpr();
    ApplicativeExprContext *applicativeExpr();
    BasePrimaryExprContext *basePrimaryExpr();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ApplicativeExprContext* applicativeExpr();
  ApplicativeExprContext* applicativeExpr(int precedence);
  class  TypeApplicativeExprContext : public antlr4::ParserRuleContext {
  public:
    TypeApplicativeExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    PrimaryExprContext *primaryExpr();
    TypeApplicativeArgsContext *typeApplicativeArgs();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TypeApplicativeExprContext* typeApplicativeExpr();

  class  TypeApplicativeArgsContext : public antlr4::ParserRuleContext {
  public:
    TypeApplicativeArgsContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *LBRACKET();
    TupleTypeArgsContext *tupleTypeArgs();
    antlr4::tree::TerminalNode *RBRACKET();
    Type_Context *type_();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TypeApplicativeArgsContext* typeApplicativeArgs();

  class  BasePrimaryExprContext : public antlr4::ParserRuleContext {
  public:
    BasePrimaryExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *AMP();
    ProjectionExprContext *projectionExpr();
    antlr4::tree::TerminalNode *INT_LITERAL();
    antlr4::tree::TerminalNode *FLOAT_LITERAL();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  BasePrimaryExprContext* basePrimaryExpr();

  class  PrimaryExprContext : public antlr4::ParserRuleContext {
  public:
    PrimaryExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *MUL();
    PrimaryExprContext *primaryExpr();
    BasePrimaryExprContext *basePrimaryExpr();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  PrimaryExprContext* primaryExpr();

  class  ProjectionExprContext : public antlr4::ParserRuleContext {
  public:
    ProjectionExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    DerefExprContext *derefExpr();
    ProjectionExprContext *projectionExpr();
    antlr4::tree::TerminalNode *DOT();
    antlr4::tree::TerminalNode *INT_LITERAL();
    antlr4::tree::TerminalNode *IDENTIFIER();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  ProjectionExprContext* projectionExpr();
  ProjectionExprContext* projectionExpr(int precedence);
  class  DerefExprContext : public antlr4::ParserRuleContext {
  public:
    DerefExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    InjectionExprContext *injectionExpr();
    antlr4::tree::TerminalNode *RUNE_LITERAL();
    antlr4::tree::TerminalNode *STRING_LITERAL();
    StructExprContext *structExpr();
    TupleExprContext *tupleExpr();
    VarExprContext *varExpr();
    antlr4::tree::TerminalNode *LPAREN();
    ExpressionContext *expression();
    antlr4::tree::TerminalNode *RPAREN();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  DerefExprContext* derefExpr();

  class  InjectionExprContext : public antlr4::ParserRuleContext {
  public:
    InjectionExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *VARIANT();
    antlr4::tree::TerminalNode *LBRACE();
    Type_Context *type_();
    LabelValueContext *labelValue();
    antlr4::tree::TerminalNode *RBRACE();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  InjectionExprContext* injectionExpr();

  class  StructExprContext : public antlr4::ParserRuleContext {
  public:
    StructExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *STRUCT();
    antlr4::tree::TerminalNode *LBRACE();
    antlr4::tree::TerminalNode *RBRACE();
    LabelValuesContext *labelValues();
    antlr4::tree::TerminalNode *COMMA();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  StructExprContext* structExpr();

  class  LabelValuesContext : public antlr4::ParserRuleContext {
  public:
    LabelValuesContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    std::vector<LabelValueContext *> labelValue();
    LabelValueContext* labelValue(size_t i);
    std::vector<antlr4::tree::TerminalNode *> COMMA();
    antlr4::tree::TerminalNode* COMMA(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  LabelValuesContext* labelValues();

  class  LabelValueContext : public antlr4::ParserRuleContext {
  public:
    LabelValueContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    IdContext *id();
    antlr4::tree::TerminalNode *ASSIGN();
    ExpressionContext *expression();
    antlr4::tree::TerminalNode *INT_LITERAL();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  LabelValueContext* labelValue();

  class  TupleExprContext : public antlr4::ParserRuleContext {
  public:
    TupleExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    antlr4::tree::TerminalNode *LPAREN();
    antlr4::tree::TerminalNode *RPAREN();
    TupleExprArgsContext *tupleExprArgs();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TupleExprContext* tupleExpr();

  class  TupleExprArgsContext : public antlr4::ParserRuleContext {
  public:
    TupleExprArgsContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    std::vector<ExpressionContext *> expression();
    ExpressionContext* expression(size_t i);
    std::vector<antlr4::tree::TerminalNode *> COMMA();
    antlr4::tree::TerminalNode* COMMA(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  TupleExprArgsContext* tupleExprArgs();

  class  VarExprContext : public antlr4::ParserRuleContext {
  public:
    VarExprContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    IdContext *id();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  VarExprContext* varExpr();

  class  IdContext : public antlr4::ParserRuleContext {
  public:
    IdContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    IdTokensContext *idTokens();
    antlr4::tree::TerminalNode *LPAREN();
    antlr4::tree::TerminalNode *RPAREN();
    antlr4::tree::TerminalNode *OR();
    antlr4::tree::TerminalNode *AND();
    antlr4::tree::TerminalNode *NE();
    antlr4::tree::TerminalNode *EQ();
    antlr4::tree::TerminalNode *GT();
    antlr4::tree::TerminalNode *GE();
    antlr4::tree::TerminalNode *LT();
    antlr4::tree::TerminalNode *LE();
    antlr4::tree::TerminalNode *PLUS();
    antlr4::tree::TerminalNode *MINUS();
    antlr4::tree::TerminalNode *MUL();
    antlr4::tree::TerminalNode *DIV();
    antlr4::tree::TerminalNode *NOT();


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  IdContext* id();

  class  IdTokensContext : public antlr4::ParserRuleContext {
  public:
    IdTokensContext(antlr4::ParserRuleContext *parent, size_t invokingState);
    virtual size_t getRuleIndex() const override;
    std::vector<antlr4::tree::TerminalNode *> IDENTIFIER();
    antlr4::tree::TerminalNode* IDENTIFIER(size_t i);
    std::vector<antlr4::tree::TerminalNode *> DOUBLE_COLON();
    antlr4::tree::TerminalNode* DOUBLE_COLON(size_t i);
    std::vector<antlr4::tree::TerminalNode *> SET();
    antlr4::tree::TerminalNode* SET(size_t i);


    virtual antlrcpp::Any accept(antlr4::tree::ParseTreeVisitor *visitor) override;
   
  };

  IdTokensContext* idTokens();


  virtual bool sempred(antlr4::RuleContext *_localctx, size_t ruleIndex, size_t predicateIndex) override;
  bool appTypeSempred(AppTypeContext *_localctx, size_t predicateIndex);
  bool logicalOrExprSempred(LogicalOrExprContext *_localctx, size_t predicateIndex);
  bool logicalAndExprSempred(LogicalAndExprContext *_localctx, size_t predicateIndex);
  bool equalityExprSempred(EqualityExprContext *_localctx, size_t predicateIndex);
  bool comparisonExprSempred(ComparisonExprContext *_localctx, size_t predicateIndex);
  bool additiveExprSempred(AdditiveExprContext *_localctx, size_t predicateIndex);
  bool multiplicativeExprSempred(MultiplicativeExprContext *_localctx, size_t predicateIndex);
  bool applicativeExprSempred(ApplicativeExprContext *_localctx, size_t predicateIndex);
  bool projectionExprSempred(ProjectionExprContext *_localctx, size_t predicateIndex);

private:
  static std::vector<antlr4::dfa::DFA> _decisionToDFA;
  static antlr4::atn::PredictionContextCache _sharedContextCache;
  static std::vector<std::string> _ruleNames;
  static std::vector<std::string> _tokenNames;

  static std::vector<std::string> _literalNames;
  static std::vector<std::string> _symbolicNames;
  static antlr4::dfa::Vocabulary _vocabulary;
  static antlr4::atn::ATN _atn;
  static std::vector<uint16_t> _serializedATN;


  struct Initializer {
    Initializer();
  };
  static Initializer _init;
};

