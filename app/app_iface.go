package app

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/ecdsa"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/gorilla/websocket"
	"github.com/mattermost/go-i18n/i18n"
	goi18n "github.com/mattermost/go-i18n/i18n"
	"github.com/mattermost/mattermost-server/v5/einterfaces"
	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/services/filesstore"
	"github.com/mattermost/mattermost-server/v5/services/httpservice"
	"github.com/mattermost/mattermost-server/v5/services/imageproxy"
	"github.com/mattermost/mattermost-server/v5/services/timezones"
	"github.com/mattermost/mattermost-server/v5/store"
	"github.com/mattermost/mattermost-server/v5/utils"
)

type AppIface interface {
	AcceptLanguage() string
	AccountMigration() einterfaces.AccountMigrationInterface
	ActivateMfa(userId, token string) *model.AppError
	AddChannelMember(userId string, channel *model.Channel, userRequestorId string, postRootId string) (*model.ChannelMember, *model.AppError)
	AddConfigListener(listener func(*model.Config, *model.Config)) string
	AddCursorIdsForPostList(originalList *model.PostList, afterPost, beforePost string, since int64, page, perPage int)
	AddDirectChannels(teamId string, user *model.User) *model.AppError
	AddLicenseListener(listener func()) string
	AddNotificationEmailToBatch(user *model.User, post *model.Post, team *model.Team) *model.AppError
	AddPublicKey(name string, key io.Reader) *model.AppError
	AddSamlIdpCertificate(fileData *multipart.FileHeader) *model.AppError
	AddSamlPrivateCertificate(fileData *multipart.FileHeader) *model.AppError
	AddSamlPublicCertificate(fileData *multipart.FileHeader) *model.AppError
	AddSessionToCache(session *model.Session)
	AddStatusCache(status *model.Status)
	AddStatusCacheSkipClusterSend(status *model.Status)
	AddTeamMember(teamId, userId string) (*model.TeamMember, *model.AppError)
	AddTeamMemberByInviteId(inviteId, userId string) (*model.TeamMember, *model.AppError)
	AddTeamMemberByToken(userId, tokenId string) (*model.TeamMember, *model.AppError)
	AddTeamMembers(teamId string, userIds []string, userRequestorId string, graceful bool) ([]*model.TeamMemberWithError, *model.AppError)
	AddUserToChannel(user *model.User, channel *model.Channel) (*model.ChannelMember, *model.AppError)
	AddUserToTeam(teamId string, userId string, userRequestorId string) (*model.Team, *model.AppError)
	AddUserToTeamByInviteId(inviteId string, userId string) (*model.Team, *model.AppError)
	AddUserToTeamByTeamId(teamId string, user *model.User) *model.AppError
	AddUserToTeamByToken(userId string, tokenId string) (*model.Team, *model.AppError)
	AllowOAuthAppAccessToUser(userId string, authRequest *model.AuthorizeRequest) (string, *model.AppError)
	AsymmetricSigningKey() *ecdsa.PrivateKey
	AttachDeviceId(sessionId string, deviceId string, expiresAt int64) *model.AppError
	AttachSessionCookies(w http.ResponseWriter, r *http.Request)
	AuthenticateUserForLogin(id, loginId, password, mfaToken string, ldapOnly bool) (*model.User, *model.AppError)
	AuthorizeOAuthUser(w http.ResponseWriter, r *http.Request, service, code, state, redirectUri string) (io.ReadCloser, string, map[string]string, *model.AppError)
	AutocompleteChannels(teamId string, term string) (*model.ChannelList, *model.AppError)
	AutocompleteChannelsForSearch(teamId string, userId string, term string) (*model.ChannelList, *model.AppError)
	AutocompleteUsersInChannel(teamId string, channelId string, term string, options *model.UserSearchOptions) (*model.UserAutocompleteInChannel, *model.AppError)
	AutocompleteUsersInTeam(teamId string, term string, options *model.UserSearchOptions) (*model.UserAutocompleteInTeam, *model.AppError)
	BroadcastStatus(status *model.Status)
	BuildPostReactions(postId string) (*[]ReactionImportData, *model.AppError)
	BuildPushNotificationMessage(contentsConfig string, post *model.Post, user *model.User, channel *model.Channel, channelName string, senderName string, explicitMention bool, channelWideMention bool, replyToThreadType string) (*model.PushNotification, *model.AppError)
	BuildSamlMetadataObject(idpMetadata []byte) (*model.SamlMetadataResponse, *model.AppError)
	BulkExport(writer io.Writer, file string, pathToEmojiDir string, dirNameToExportEmoji string) *model.AppError
	BulkImport(fileReader io.Reader, dryRun bool, workers int) (*model.AppError, int)
	CancelJob(jobId string) *model.AppError
	ChannelMembersMinusGroupMembers(channelID string, groupIDs []string, page, perPage int) ([]*model.UserWithGroups, int64, *model.AppError)
	ChannelMembersToAdd(since int64, channelID *string) ([]*model.UserChannelIDPair, *model.AppError)
	ChannelMembersToRemove(teamID *string) ([]*model.ChannelMember, *model.AppError)
	CheckForClientSideCert(r *http.Request) (string, string, string)
	CheckPasswordAndAllCriteria(user *model.User, password string, mfaToken string) *model.AppError
	CheckRolesExist(roleNames []string) *model.AppError
	CheckUserAllAuthenticationCriteria(user *model.User, mfaToken string) *model.AppError
	CheckUserMfa(user *model.User, token string) *model.AppError
	CheckUserPostflightAuthenticationCriteria(user *model.User) *model.AppError
	CheckUserPreflightAuthenticationCriteria(user *model.User, mfaToken string) *model.AppError
	ClearChannelMembersCache(channelID string)
	ClearPushNotification(currentSessionId, userId, channelId string)
	ClearPushNotificationSync(currentSessionId, userId, channelId string) *model.AppError
	ClearSessionCacheForAllUsers()
	ClearSessionCacheForAllUsersSkipClusterSend()
	ClearSessionCacheForUser(userId string)
	ClearSessionCacheForUserSkipClusterSend(userId string)
	ClearTeamMembersCache(teamID string)
	ClientConfig() map[string]string
	ClientConfigHash() string
	ClientConfigWithComputed() map[string]string
	ClientLicense() map[string]string
	Cluster() einterfaces.ClusterInterface
	CompareAndDeletePluginKey(pluginId string, key string, oldValue []byte) (bool, *model.AppError)
	CompareAndSetPluginKey(pluginId string, key string, oldValue, newValue []byte) (bool, *model.AppError)
	CompleteOAuth(service string, body io.ReadCloser, teamId string, props map[string]string) (*model.User, *model.AppError)
	CompleteSwitchWithOAuth(service string, userData io.Reader, email string) (*model.User, *model.AppError)
	Compliance() einterfaces.ComplianceInterface
	Config() *model.Config
	Context() context.Context
	ConvertUserToBot(user *model.User) (*model.Bot, *model.AppError)
	CopyFileInfos(userId string, fileIds []string) ([]string, *model.AppError)
	CreateBasicUser(client *model.Client4) *model.AppError
	CreateBot(bot *model.Bot) (*model.Bot, *model.AppError)
	CreateChannel(channel *model.Channel, addMember bool) (*model.Channel, *model.AppError)
	CreateChannelWithUser(channel *model.Channel, userId string) (*model.Channel, *model.AppError)
	CreateCommand(cmd *model.Command) (*model.Command, *model.AppError)
	CreateCommandPost(post *model.Post, teamId string, response *model.CommandResponse, skipSlackParsing bool) (*model.Post, *model.AppError)
	CreateCommandWebhook(commandId string, args *model.CommandArgs) (*model.CommandWebhook, *model.AppError)
	CreateDefaultChannels(teamID string) ([]*model.Channel, *model.AppError)
	CreateDefaultMemberships(since int64) error
	CreateEmoji(sessionUserId string, emoji *model.Emoji, multiPartImageData *multipart.Form) (*model.Emoji, *model.AppError)
	CreateGroup(group *model.Group) (*model.Group, *model.AppError)
	CreateGroupChannel(userIds []string, creatorId string) (*model.Channel, *model.AppError)
	CreateGuest(user *model.User) (*model.User, *model.AppError)
	CreateIncomingWebhookForChannel(creatorId string, channel *model.Channel, hook *model.IncomingWebhook) (*model.IncomingWebhook, *model.AppError)
	CreateJob(job *model.Job) (*model.Job, *model.AppError)
	CreateOAuthApp(app *model.OAuthApp) (*model.OAuthApp, *model.AppError)
	CreateOAuthStateToken(extra string) (*model.Token, *model.AppError)
	CreateOAuthUser(service string, userData io.Reader, teamId string) (*model.User, *model.AppError)
	CreateOutgoingWebhook(hook *model.OutgoingWebhook) (*model.OutgoingWebhook, *model.AppError)
	CreatePasswordRecoveryToken(userId, email string) (*model.Token, *model.AppError)
	CreatePost(post *model.Post, channel *model.Channel, triggerWebhooks bool) (*model.Post, *model.AppError)
	CreatePostAsUser(post *model.Post, currentSessionId string) (*model.Post, *model.AppError)
	CreatePostMissingChannel(post *model.Post, triggerWebhooks bool) (*model.Post, *model.AppError)
	CreatePushNotificationsHub()
	CreateRole(role *model.Role) (*model.Role, *model.AppError)
	CreateScheme(scheme *model.Scheme) (*model.Scheme, *model.AppError)
	CreateSession(session *model.Session) (*model.Session, *model.AppError)
	CreateTeam(team *model.Team) (*model.Team, *model.AppError)
	CreateTeamWithUser(team *model.Team, userId string) (*model.Team, *model.AppError)
	CreateTermsOfService(text, userId string) (*model.TermsOfService, *model.AppError)
	CreateUser(user *model.User) (*model.User, *model.AppError)
	CreateUserAccessToken(token *model.UserAccessToken) (*model.UserAccessToken, *model.AppError)
	CreateUserAsAdmin(user *model.User) (*model.User, *model.AppError)
	CreateUserFromSignup(user *model.User) (*model.User, *model.AppError)
	CreateUserWithInviteId(user *model.User, inviteId string) (*model.User, *model.AppError)
	CreateUserWithToken(user *model.User, token *model.Token) (*model.User, *model.AppError)
	CreateVerifyEmailToken(userId string, newEmail string) (*model.Token, *model.AppError)
	CreateWebhookPost(userId string, channel *model.Channel, text, overrideUsername, overrideIconUrl, overrideIconEmoji string, props model.StringInterface, postType string, postRootId string) (*model.Post, *model.AppError)
	DataRetention() einterfaces.DataRetentionInterface
	DeactivateGuests() *model.AppError
	DeactivateMfa(userId string) *model.AppError
	DeauthorizeOAuthAppForUser(userId, appId string) *model.AppError
	DefaultChannelNames() []string
	DeleteAllExpiredPluginKeys() *model.AppError
	DeleteAllKeysForPlugin(pluginId string) *model.AppError
	DeleteBotIconImage(botUserId string) *model.AppError
	DeleteBrandImage() *model.AppError
	DeleteChannel(channel *model.Channel, userId string) *model.AppError
	DeleteCommand(commandId string) *model.AppError
	DeleteEmoji(emoji *model.Emoji) *model.AppError
	DeleteEphemeralPost(userId, postId string)
	DeleteFlaggedPosts(postId string)
	DeleteGroup(groupID string) (*model.Group, *model.AppError)
	DeleteGroupConstrainedMemberships() error
	DeleteGroupMember(groupID string, userID string) (*model.GroupMember, *model.AppError)
	DeleteGroupSyncable(groupID string, syncableID string, syncableType model.GroupSyncableType) (*model.GroupSyncable, *model.AppError)
	DeleteIncomingWebhook(hookId string) *model.AppError
	DeleteOAuthApp(appId string) *model.AppError
	DeleteOutgoingWebhook(hookId string) *model.AppError
	DeletePluginKey(pluginId string, key string) *model.AppError
	DeletePost(postId, deleteByID string) (*model.Post, *model.AppError)
	DeletePostFiles(post *model.Post)
	DeletePreferences(userId string, preferences model.Preferences) *model.AppError
	DeletePublicKey(name string) *model.AppError
	DeleteReactionForPost(reaction *model.Reaction) *model.AppError
	DeleteScheme(schemeId string) (*model.Scheme, *model.AppError)
	DeleteToken(token *model.Token) *model.AppError
	DemoteUserToGuest(user *model.User) *model.AppError
	DiagnosticId() string
	DisableAutoResponder(userId string, asAdmin bool) *model.AppError
	DisablePlugin(id string) *model.AppError
	DisableUserAccessToken(token *model.UserAccessToken) *model.AppError
	DoActionRequest(rawURL string, body []byte) (*http.Response, *model.AppError)
	DoAdvancedPermissionsMigration()
	DoAppMigrations()
	DoEmojisPermissionsMigration()
	DoGuestRolesCreationMigration()
	DoLocalRequest(rawURL string, body []byte) (*http.Response, *model.AppError)
	DoLogin(w http.ResponseWriter, r *http.Request, user *model.User, deviceId string) *model.AppError
	DoPermissionsMigrations() *model.AppError
	DoPostAction(postId, actionId, userId, selectedOption string) (string, *model.AppError)
	DoPostActionWithCookie(postId, actionId, userId, selectedOption string, cookie *model.PostActionCookie) (string, *model.AppError)
	DoUploadFile(now time.Time, rawTeamId string, rawChannelId string, rawUserId string, rawFilename string, data []byte) (*model.FileInfo, *model.AppError)
	DoUploadFileExpectModification(now time.Time, rawTeamId string, rawChannelId string, rawUserId string, rawFilename string, data []byte) (*model.FileInfo, []byte, *model.AppError)
	DoubleCheckPassword(user *model.User, password string) *model.AppError
	DownloadFromURL(downloadURL string) ([]byte, error)
	Elasticsearch() einterfaces.ElasticsearchInterface
	EnablePlugin(id string) *model.AppError
	EnableUserAccessToken(token *model.UserAccessToken) *model.AppError
	EnsureDiagnosticId()
	EnvironmentConfig() map[string]interface{}
	ExecuteCommand(args *model.CommandArgs) (*model.CommandResponse, *model.AppError)
	ExportAllChannels(writer io.Writer) *model.AppError
	ExportAllDirectChannels(writer io.Writer) *model.AppError
	ExportAllDirectPosts(writer io.Writer) *model.AppError
	ExportAllPosts(writer io.Writer) *model.AppError
	ExportAllTeams(writer io.Writer) *model.AppError
	ExportAllUsers(writer io.Writer) *model.AppError
	ExportCustomEmoji(writer io.Writer, file string, pathToEmojiDir string, dirNameToExportEmoji string) *model.AppError
	ExportPermissions(w io.Writer) error
	ExportVersion(writer io.Writer) *model.AppError
	ExportWriteLine(writer io.Writer, line *LineImportData) *model.AppError
	FetchSamlMetadataFromIdp(url string) ([]byte, *model.AppError)
	FileBackend() (filesstore.FileBackend, *model.AppError)
	FileExists(path string) (bool, *model.AppError)
	FileReader(path string) (filesstore.ReadCloseSeeker, *model.AppError)
	FillInChannelProps(channel *model.Channel) *model.AppError
	FillInChannelsProps(channelList *model.ChannelList) *model.AppError
	FillInPostProps(post *model.Post, channel *model.Channel) *model.AppError
	FilterNonGroupChannelMembers(userIDs []string, channel *model.Channel) ([]string, error)
	FilterNonGroupTeamMembers(userIDs []string, team *model.Team) ([]string, error)
	FindTeamByName(name string) bool
	GenerateMfaSecret(userId string) (*model.MfaSecret, *model.AppError)
	GeneratePublicLink(siteURL string, info *model.FileInfo) string
	GetActivePluginManifests() ([]*model.Manifest, *model.AppError)
	GetAllChannels(page, perPage int, opts model.ChannelSearchOpts) (*model.ChannelListWithTeamData, *model.AppError)
	GetAllChannelsCount(opts model.ChannelSearchOpts) (int64, *model.AppError)
	GetAllLdapGroupsPage(page int, perPage int, opts model.LdapGroupSearchOpts) ([]*model.Group, int, *model.AppError)
	GetAllPrivateTeams() ([]*model.Team, *model.AppError)
	GetAllPrivateTeamsPage(offset int, limit int) ([]*model.Team, *model.AppError)
	GetAllPrivateTeamsPageWithCount(offset int, limit int) (*model.TeamsWithCount, *model.AppError)
	GetAllPublicTeams() ([]*model.Team, *model.AppError)
	GetAllPublicTeamsPage(offset int, limit int) ([]*model.Team, *model.AppError)
	GetAllPublicTeamsPageWithCount(offset int, limit int) (*model.TeamsWithCount, *model.AppError)
	GetAllRoles() ([]*model.Role, *model.AppError)
	GetAllStatuses() map[string]*model.Status
	GetAllTeams() ([]*model.Team, *model.AppError)
	GetAllTeamsPage(offset int, limit int) ([]*model.Team, *model.AppError)
	GetAllTeamsPageWithCount(offset int, limit int) (*model.TeamsWithCount, *model.AppError)
	GetAnalytics(name string, teamId string) (model.AnalyticsRows, *model.AppError)
	GetAudits(userId string, limit int) (model.Audits, *model.AppError)
	GetAuditsPage(userId string, page int, perPage int) (model.Audits, *model.AppError)
	GetAuthorizationCode(w http.ResponseWriter, r *http.Request, service string, props map[string]string, loginHint string) (string, *model.AppError)
	GetAuthorizedAppsForUser(userId string, page, perPage int) ([]*model.OAuthApp, *model.AppError)
	GetBot(botUserId string, includeDeleted bool) (*model.Bot, *model.AppError)
	GetBotIconImage(botUserId string) ([]byte, *model.AppError)
	GetBots(options *model.BotGetOptions) (model.BotList, *model.AppError)
	GetBrandImage() ([]byte, *model.AppError)
	GetBulkReactionsForPosts(postIds []string) (map[string][]*model.Reaction, *model.AppError)
	GetChannel(channelId string) (*model.Channel, *model.AppError)
	GetChannelByName(channelName, teamId string, includeDeleted bool) (*model.Channel, *model.AppError)
	GetChannelByNameForTeamName(channelName, teamName string, includeDeleted bool) (*model.Channel, *model.AppError)
	GetChannelCounts(teamId string, userId string) (*model.ChannelCounts, *model.AppError)
	GetChannelGroupUsers(channelID string) ([]*model.User, *model.AppError)
	GetChannelGuestCount(channelId string) (int64, *model.AppError)
	GetChannelMember(channelId string, userId string) (*model.ChannelMember, *model.AppError)
	GetChannelMemberCount(channelId string) (int64, *model.AppError)
	GetChannelMembersByIds(channelId string, userIds []string) (*model.ChannelMembers, *model.AppError)
	GetChannelMembersForUser(teamId string, userId string) (*model.ChannelMembers, *model.AppError)
	GetChannelMembersForUserWithPagination(teamId, userId string, page, perPage int) ([]*model.ChannelMember, *model.AppError)
	GetChannelMembersPage(channelId string, page, perPage int) (*model.ChannelMembers, *model.AppError)
	GetChannelMembersTimezones(channelId string) ([]string, *model.AppError)
	GetChannelPinnedPostCount(channelId string) (int64, *model.AppError)
	GetChannelUnread(channelId, userId string) (*model.ChannelUnread, *model.AppError)
	GetChannelsByNames(channelNames []string, teamId string) ([]*model.Channel, *model.AppError)
	GetChannelsForScheme(scheme *model.Scheme, offset int, limit int) (model.ChannelList, *model.AppError)
	GetChannelsForSchemePage(scheme *model.Scheme, page int, perPage int) (model.ChannelList, *model.AppError)
	GetChannelsForUser(teamId string, userId string, includeDeleted bool) (*model.ChannelList, *model.AppError)
	GetChannelsUserNotIn(teamId string, userId string, offset int, limit int) (*model.ChannelList, *model.AppError)
	GetClusterId() string
	GetClusterPluginStatuses() (model.PluginStatuses, *model.AppError)
	GetClusterStatus() []*model.ClusterInfo
	GetCommand(commandId string) (*model.Command, *model.AppError)
	GetComplianceFile(job *model.Compliance) ([]byte, *model.AppError)
	GetComplianceReport(reportId string) (*model.Compliance, *model.AppError)
	GetComplianceReports(page, perPage int) (model.Compliances, *model.AppError)
	GetConfigFile(name string) ([]byte, error)
	GetCookieDomain() string
	GetDataRetentionPolicy() (*model.DataRetentionPolicy, *model.AppError)
	GetDefaultProfileImage(user *model.User) ([]byte, *model.AppError)
	GetDeletedChannels(teamId string, offset int, limit int, userId string) (*model.ChannelList, *model.AppError)
	GetEmoji(emojiId string) (*model.Emoji, *model.AppError)
	GetEmojiByName(emojiName string) (*model.Emoji, *model.AppError)
	GetEmojiImage(emojiId string) ([]byte, string, *model.AppError)
	GetEmojiList(page, perPage int, sort string) ([]*model.Emoji, *model.AppError)
	GetEmojiStaticUrl(emojiName string) (string, *model.AppError)
	GetEnvironmentConfig() map[string]interface{}
	GetFile(fileId string) ([]byte, *model.AppError)
	GetFileInfo(fileId string) (*model.FileInfo, *model.AppError)
	GetFileInfosForPost(postId string, fromMaster bool) ([]*model.FileInfo, *model.AppError)
	GetFileInfosForPostWithMigration(postId string) ([]*model.FileInfo, *model.AppError)
	GetFlaggedPosts(userId string, offset int, limit int) (*model.PostList, *model.AppError)
	GetFlaggedPostsForChannel(userId, channelId string, offset int, limit int) (*model.PostList, *model.AppError)
	GetFlaggedPostsForTeam(userId, teamId string, offset int, limit int) (*model.PostList, *model.AppError)
	GetGroup(id string) (*model.Group, *model.AppError)
	GetGroupByName(name string) (*model.Group, *model.AppError)
	GetGroupByRemoteID(remoteID string, groupSource model.GroupSource) (*model.Group, *model.AppError)
	GetGroupChannel(userIds []string) (*model.Channel, *model.AppError)
	GetGroupMemberUsers(groupID string) ([]*model.User, *model.AppError)
	GetGroupMemberUsersPage(groupID string, page int, perPage int) ([]*model.User, int, *model.AppError)
	GetGroupSyncable(groupID string, syncableID string, syncableType model.GroupSyncableType) (*model.GroupSyncable, *model.AppError)
	GetGroupSyncables(groupID string, syncableType model.GroupSyncableType) ([]*model.GroupSyncable, *model.AppError)
	GetGroups(page, perPage int, opts model.GroupSearchOpts) ([]*model.Group, *model.AppError)
	GetGroupsByChannel(channelId string, opts model.GroupSearchOpts) ([]*model.GroupWithSchemeAdmin, int, *model.AppError)
	GetGroupsByIDs(groupIDs []string) ([]*model.Group, *model.AppError)
	GetGroupsBySource(groupSource model.GroupSource) ([]*model.Group, *model.AppError)
	GetGroupsByTeam(teamId string, opts model.GroupSearchOpts) ([]*model.GroupWithSchemeAdmin, int, *model.AppError)
	GetGroupsByUserId(userId string) ([]*model.Group, *model.AppError)
	GetHubForUserId(userId string) *Hub
	GetIncomingWebhook(hookId string) (*model.IncomingWebhook, *model.AppError)
	GetIncomingWebhooksForTeamPage(teamId string, page, perPage int) ([]*model.IncomingWebhook, *model.AppError)
	GetIncomingWebhooksForTeamPageByUser(teamId string, userId string, page, perPage int) ([]*model.IncomingWebhook, *model.AppError)
	GetIncomingWebhooksPage(page, perPage int) ([]*model.IncomingWebhook, *model.AppError)
	GetIncomingWebhooksPageByUser(userId string, page, perPage int) ([]*model.IncomingWebhook, *model.AppError)
	GetJob(id string) (*model.Job, *model.AppError)
	GetJobs(offset int, limit int) ([]*model.Job, *model.AppError)
	GetJobsByType(jobType string, offset int, limit int) ([]*model.Job, *model.AppError)
	GetJobsByTypePage(jobType string, page int, perPage int) ([]*model.Job, *model.AppError)
	GetJobsPage(page int, perPage int) ([]*model.Job, *model.AppError)
	GetLatestTermsOfService() (*model.TermsOfService, *model.AppError)
	GetLdapGroup(ldapGroupID string) (*model.Group, *model.AppError)
	GetLogs(page, perPage int) ([]string, *model.AppError)
	GetLogsSkipSend(page, perPage int) ([]string, *model.AppError)
	GetMarketplacePlugins(filter *model.MarketplacePluginFilter) ([]*model.MarketplacePlugin, *model.AppError)
	GetMessageForNotification(post *model.Post, translateFunc i18n.TranslateFunc) string
	GetMultipleEmojiByName(names []string) ([]*model.Emoji, *model.AppError)
	GetNewUsersForTeamPage(teamId string, page, perPage int, asAdmin bool, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, *model.AppError)
	GetNextPostIdFromPostList(postList *model.PostList) string
	GetNotificationNameFormat(user *model.User) string
	GetNumberOfChannelsOnTeam(teamId string) (int, *model.AppError)
	GetOAuthAccessTokenForCodeFlow(clientId, grantType, redirectUri, code, secret, refreshToken string) (*model.AccessResponse, *model.AppError)
	GetOAuthAccessTokenForImplicitFlow(userId string, authRequest *model.AuthorizeRequest) (*model.Session, *model.AppError)
	GetOAuthApp(appId string) (*model.OAuthApp, *model.AppError)
	GetOAuthApps(page, perPage int) ([]*model.OAuthApp, *model.AppError)
	GetOAuthAppsByCreator(userId string, page, perPage int) ([]*model.OAuthApp, *model.AppError)
	GetOAuthCodeRedirect(userId string, authRequest *model.AuthorizeRequest) (string, *model.AppError)
	GetOAuthImplicitRedirect(userId string, authRequest *model.AuthorizeRequest) (string, *model.AppError)
	GetOAuthLoginEndpoint(w http.ResponseWriter, r *http.Request, service, teamId, action, redirectTo, loginHint string) (string, *model.AppError)
	GetOAuthSignupEndpoint(w http.ResponseWriter, r *http.Request, service, teamId string) (string, *model.AppError)
	GetOAuthStateToken(token string) (*model.Token, *model.AppError)
	GetOpenGraphMetadata(requestURL string) *opengraph.OpenGraph
	GetOrCreateDirectChannel(userId, otherUserId string) (*model.Channel, *model.AppError)
	GetOutgoingWebhook(hookId string) (*model.OutgoingWebhook, *model.AppError)
	GetOutgoingWebhooksForChannelPageByUser(channelId string, userId string, page, perPage int) ([]*model.OutgoingWebhook, *model.AppError)
	GetOutgoingWebhooksForTeamPage(teamId string, page, perPage int) ([]*model.OutgoingWebhook, *model.AppError)
	GetOutgoingWebhooksForTeamPageByUser(teamId string, userId string, page, perPage int) ([]*model.OutgoingWebhook, *model.AppError)
	GetOutgoingWebhooksPage(page, perPage int) ([]*model.OutgoingWebhook, *model.AppError)
	GetOutgoingWebhooksPageByUser(userId string, page, perPage int) ([]*model.OutgoingWebhook, *model.AppError)
	GetPasswordRecoveryToken(token string) (*model.Token, *model.AppError)
	GetPermalinkPost(postId string, userId string) (*model.PostList, *model.AppError)
	GetPinnedPosts(channelId string) (*model.PostList, *model.AppError)
	GetPluginKey(pluginId string, key string) ([]byte, *model.AppError)
	GetPluginPublicKeyFiles() ([]string, *model.AppError)
	GetPluginStatus(id string) (*model.PluginStatus, *model.AppError)
	GetPluginStatuses() (model.PluginStatuses, *model.AppError)
	GetPlugins() (*model.PluginsResponse, *model.AppError)
	GetPluginsEnvironment() *plugin.Environment
	GetPostAfterTime(channelId string, time int64) (*model.Post, *model.AppError)
	GetPostIdAfterTime(channelId string, time int64) (string, *model.AppError)
	GetPostIdBeforeTime(channelId string, time int64) (string, *model.AppError)
	GetPostThread(postId string) (*model.PostList, *model.AppError)
	GetPosts(channelId string, offset int, limit int) (*model.PostList, *model.AppError)
	GetPostsAfterPost(channelId, postId string, page, perPage int) (*model.PostList, *model.AppError)
	GetPostsAroundPost(postId, channelId string, offset, limit int, before bool) (*model.PostList, *model.AppError)
	GetPostsBeforePost(channelId, postId string, page, perPage int) (*model.PostList, *model.AppError)
	GetPostsEtag(channelId string) string
	GetPostsForChannelAroundLastUnread(channelId, userId string, limitBefore, limitAfter int) (*model.PostList, *model.AppError)
	GetPostsPage(channelId string, page int, perPage int) (*model.PostList, *model.AppError)
	GetPostsSince(channelId string, time int64) (*model.PostList, *model.AppError)
	GetPreferenceByCategoryAndNameForUser(userId string, category string, preferenceName string) (*model.Preference, *model.AppError)
	GetPreferenceByCategoryForUser(userId string, category string) (model.Preferences, *model.AppError)
	GetPreferencesForUser(userId string) (model.Preferences, *model.AppError)
	GetPrevPostIdFromPostList(postList *model.PostList) string
	GetProfileImage(user *model.User) ([]byte, bool, *model.AppError)
	GetPublicChannelsByIdsForTeam(teamId string, channelIds []string) (*model.ChannelList, *model.AppError)
	GetPublicChannelsForTeam(teamId string, offset int, limit int) (*model.ChannelList, *model.AppError)
	GetPublicKey(name string) ([]byte, *model.AppError)
	GetReactionsForPost(postId string) ([]*model.Reaction, *model.AppError)
	GetRecentlyActiveUsersForTeam(teamId string) (map[string]*model.User, *model.AppError)
	GetRecentlyActiveUsersForTeamPage(teamId string, page, perPage int, asAdmin bool, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, *model.AppError)
	GetRole(id string) (*model.Role, *model.AppError)
	GetRoleByName(name string) (*model.Role, *model.AppError)
	GetRolesByNames(names []string) ([]*model.Role, *model.AppError)
	GetSamlCertificateStatus() *model.SamlCertificateStatus
	GetSamlMetadata() (string, *model.AppError)
	GetSamlMetadataFromIdp(idpMetadataUrl string) (*model.SamlMetadataResponse, *model.AppError)
	GetSanitizeOptions(asAdmin bool) map[string]bool
	GetSanitizedClientLicense() map[string]string
	GetSanitizedConfig() *model.Config
	GetScheme(id string) (*model.Scheme, *model.AppError)
	GetSchemeByName(name string) (*model.Scheme, *model.AppError)
	GetSchemeRolesForChannel(channelId string) (string, string, string, *model.AppError)
	GetSchemeRolesForTeam(teamId string) (string, string, string, *model.AppError)
	GetSchemes(scope string, offset int, limit int) ([]*model.Scheme, *model.AppError)
	GetSchemesPage(scope string, page int, perPage int) ([]*model.Scheme, *model.AppError)
	GetSession(token string) (*model.Session, *model.AppError)
	GetSessionById(sessionId string) (*model.Session, *model.AppError)
	GetSessions(userId string) ([]*model.Session, *model.AppError)
	GetSinglePost(postId string) (*model.Post, *model.AppError)
	GetSiteURL() string
	GetStatus(userId string) (*model.Status, *model.AppError)
	GetStatusFromCache(userId string) *model.Status
	GetStatusesByIds(userIds []string) (map[string]interface{}, *model.AppError)
	GetT() goi18n.TranslateFunc
	GetTeam(teamId string) (*model.Team, *model.AppError)
	GetTeamByInviteId(inviteId string) (*model.Team, *model.AppError)
	GetTeamByName(name string) (*model.Team, *model.AppError)
	GetTeamGroupUsers(teamID string) ([]*model.User, *model.AppError)
	GetTeamIcon(team *model.Team) ([]byte, *model.AppError)
	GetTeamIdFromQuery(query url.Values) (string, *model.AppError)
	GetTeamMember(teamId, userId string) (*model.TeamMember, *model.AppError)
	GetTeamMembers(teamId string, offset int, limit int, restrictions *model.ViewUsersRestrictions) ([]*model.TeamMember, *model.AppError)
	GetTeamMembersByIds(teamId string, userIds []string, restrictions *model.ViewUsersRestrictions) ([]*model.TeamMember, *model.AppError)
	GetTeamMembersForUser(userId string) ([]*model.TeamMember, *model.AppError)
	GetTeamMembersForUserWithPagination(userId string, page, perPage int) ([]*model.TeamMember, *model.AppError)
	GetTeamStats(teamId string, restrictions *model.ViewUsersRestrictions) (*model.TeamStats, *model.AppError)
	GetTeamUnread(teamId, userId string) (*model.TeamUnread, *model.AppError)
	GetTeamsForScheme(scheme *model.Scheme, offset int, limit int) ([]*model.Team, *model.AppError)
	GetTeamsForSchemePage(scheme *model.Scheme, page int, perPage int) ([]*model.Team, *model.AppError)
	GetTeamsForUser(userId string) ([]*model.Team, *model.AppError)
	GetTeamsUnreadForUser(excludeTeamId string, userId string) ([]*model.TeamUnread, *model.AppError)
	GetTermsOfService(id string) (*model.TermsOfService, *model.AppError)
	GetTotalUsersStats(viewRestrictions *model.ViewUsersRestrictions) (*model.UsersStats, *model.AppError)
	GetUser(userId string) (*model.User, *model.AppError)
	GetUserAccessToken(tokenId string, sanitize bool) (*model.UserAccessToken, *model.AppError)
	GetUserAccessTokens(page, perPage int) ([]*model.UserAccessToken, *model.AppError)
	GetUserAccessTokensForUser(userId string, page, perPage int) ([]*model.UserAccessToken, *model.AppError)
	GetUserByAuth(authData *string, authService string) (*model.User, *model.AppError)
	GetUserByEmail(email string) (*model.User, *model.AppError)
	GetUserByUsername(username string) (*model.User, *model.AppError)
	GetUserForLogin(id, loginId string) (*model.User, *model.AppError)
	GetUserStatusesByIds(userIds []string) ([]*model.Status, *model.AppError)
	GetUserTermsOfService(userId string) (*model.UserTermsOfService, *model.AppError)
	GetUsers(options *model.UserGetOptions) ([]*model.User, *model.AppError)
	GetUsersByGroupChannelIds(channelIds []string, asAdmin bool) (map[string][]*model.User, *model.AppError)
	GetUsersByIds(userIds []string, options *store.UserGetByIdsOpts) ([]*model.User, *model.AppError)
	GetUsersByUsernames(usernames []string, asAdmin bool, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, *model.AppError)
	GetUsersEtag(restrictionsHash string) string
	GetUsersInChannel(channelId string, offset int, limit int) ([]*model.User, *model.AppError)
	GetUsersInChannelByStatus(channelId string, offset int, limit int) ([]*model.User, *model.AppError)
	GetUsersInChannelMap(channelId string, offset int, limit int, asAdmin bool) (map[string]*model.User, *model.AppError)
	GetUsersInChannelPage(channelId string, page int, perPage int, asAdmin bool) ([]*model.User, *model.AppError)
	GetUsersInChannelPageByStatus(channelId string, page int, perPage int, asAdmin bool) ([]*model.User, *model.AppError)
	GetUsersInTeam(options *model.UserGetOptions) ([]*model.User, *model.AppError)
	GetUsersInTeamEtag(teamId string, restrictionsHash string) string
	GetUsersInTeamPage(options *model.UserGetOptions, asAdmin bool) ([]*model.User, *model.AppError)
	GetUsersNotInChannel(teamId string, channelId string, groupConstrained bool, offset int, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, *model.AppError)
	GetUsersNotInChannelMap(teamId string, channelId string, groupConstrained bool, offset int, limit int, asAdmin bool, viewRestrictions *model.ViewUsersRestrictions) (map[string]*model.User, *model.AppError)
	GetUsersNotInChannelPage(teamId string, channelId string, groupConstrained bool, page int, perPage int, asAdmin bool, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, *model.AppError)
	GetUsersNotInTeam(teamId string, groupConstrained bool, offset int, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, *model.AppError)
	GetUsersNotInTeamEtag(teamId string, restrictionsHash string) string
	GetUsersNotInTeamPage(teamId string, groupConstrained bool, page int, perPage int, asAdmin bool, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, *model.AppError)
	GetUsersPage(options *model.UserGetOptions, asAdmin bool) ([]*model.User, *model.AppError)
	GetUsersWithoutTeam(options *model.UserGetOptions) ([]*model.User, *model.AppError)
	GetUsersWithoutTeamPage(options *model.UserGetOptions, asAdmin bool) ([]*model.User, *model.AppError)
	GetVerifyEmailToken(token string) (*model.Token, *model.AppError)
	GetViewUsersRestrictions(userId string) (*model.ViewUsersRestrictions, *model.AppError)
	GetViewUsersRestrictionsForTeam(userId string, teamId string) ([]string, *model.AppError)
	HTMLTemplates() *template.Template
	HTTPService() httpservice.HTTPService
	Handle404(w http.ResponseWriter, r *http.Request)
	HandleCommandResponse(command *model.Command, args *model.CommandArgs, response *model.CommandResponse, builtIn bool) (*model.CommandResponse, *model.AppError)
	HandleCommandResponsePost(command *model.Command, args *model.CommandArgs, response *model.CommandResponse, builtIn bool) (*model.Post, *model.AppError)
	HandleCommandWebhook(hookId string, response *model.CommandResponse) *model.AppError
	HandleImages(previewPathList []string, thumbnailPathList []string, fileData [][]byte)
	HandleIncomingWebhook(hookId string, req *model.IncomingWebhookRequest) *model.AppError
	HandleMessageExportConfig(cfg *model.Config, appCfg *model.Config)
	HasPermissionTo(askingUserId string, permission *model.Permission) bool
	HasPermissionToChannel(askingUserId string, channelId string, permission *model.Permission) bool
	HasPermissionToChannelByPost(askingUserId string, postId string, permission *model.Permission) bool
	HasPermissionToTeam(askingUserId string, teamId string, permission *model.Permission) bool
	HasPermissionToUser(askingUserId string, userId string) bool
	HubRegister(webConn *WebConn)
	HubStart()
	HubStop()
	HubUnregister(webConn *WebConn)
	ImageProxy() *imageproxy.ImageProxy
	ImageProxyAdder() func(string) string
	ImageProxyRemover() func(string) string
	ImportPermissions(jsonl io.Reader) error
	InitPlugins(pluginDir, webappPluginDir string)
	InitPostMetadata()
	InstallMarketplacePlugin(request *model.InstallMarketplacePluginRequest) (*model.Manifest, *model.AppError)
	InstallPlugin(pluginFile io.ReadSeeker, replace bool) (*model.Manifest, *model.AppError)
	InstallPluginFromData(data model.PluginEventData)
	InstallPluginWithSignature(pluginFile, signature io.ReadSeeker) (*model.Manifest, *model.AppError)
	InvalidateAllCaches() *model.AppError
	InvalidateAllCachesSkipSend()
	InvalidateAllEmailInvites() *model.AppError
	InvalidateCacheForChannel(channel *model.Channel)
	InvalidateCacheForChannelByNameSkipClusterSend(teamId, name string)
	InvalidateCacheForChannelMembers(channelId string)
	InvalidateCacheForChannelMembersNotifyProps(channelId string)
	InvalidateCacheForChannelMembersNotifyPropsSkipClusterSend(channelId string)
	InvalidateCacheForChannelPosts(channelId string)
	InvalidateCacheForUser(userId string)
	InvalidateCacheForUserSkipClusterSend(userId string)
	InvalidateCacheForUserTeams(userId string)
	InvalidateCacheForUserTeamsSkipClusterSend(userId string)
	InvalidateCacheForWebhook(webhookId string)
	InvalidateWebConnSessionCacheForUser(userId string)
	InviteGuestsToChannels(teamId string, guestsInvite *model.GuestsInvite, senderId string) *model.AppError
	InviteNewUsersToTeam(emailList []string, teamId, senderId string) *model.AppError
	IpAddress() string
	IsESAutocompletionEnabled() bool
	IsESIndexingEnabled() bool
	IsESSearchEnabled() bool
	IsFirstUserAccount() bool
	IsLeader() bool
	IsPasswordValid(password string) *model.AppError
	IsPhase2MigrationCompleted() *model.AppError
	IsUserAway(lastActivityAt int64) bool
	IsUserSignUpAllowed() *model.AppError
	IsUsernameTaken(name string) bool
	JoinChannel(channel *model.Channel, userId string) *model.AppError
	JoinDefaultChannels(teamId string, user *model.User, shouldBeAdmin bool, userRequestorId string) *model.AppError
	JoinUserToTeam(team *model.Team, user *model.User, userRequestorId string) *model.AppError
	Ldap() einterfaces.LdapInterface
	LeaveChannel(channelId string, userId string) *model.AppError
	LeaveTeam(team *model.Team, user *model.User, requestorId string) *model.AppError
	License() *model.License
	LimitedClientConfig() map[string]string
	LimitedClientConfigWithComputed() map[string]string
	ListAllCommands(teamId string, T goi18n.TranslateFunc) ([]*model.Command, *model.AppError)
	ListAutocompleteCommands(teamId string, T goi18n.TranslateFunc) ([]*model.Command, *model.AppError)
	ListDirectory(path string) ([]string, *model.AppError)
	ListPluginKeys(pluginId string, page, perPage int) ([]string, *model.AppError)
	ListTeamCommands(teamId string) ([]*model.Command, *model.AppError)
	LoadLicense()
	Log() *mlog.Logger
	LoginByOAuth(service string, userData io.Reader, teamId string) (*model.User, *model.AppError)
	MakePermissionError(permission *model.Permission) *model.AppError
	MarkChannelAsUnreadFromPost(postID string, userID string) (*model.ChannelUnreadAt, *model.AppError)
	MarkChannelsAsViewed(channelIds []string, userId string, currentSessionId string) (map[string]int64, *model.AppError)
	MaxPostSize() int
	MessageExport() einterfaces.MessageExportInterface
	Metrics() einterfaces.MetricsInterface
	MigrateFilenamesToFileInfos(post *model.Post) []*model.FileInfo
	MoveChannel(team *model.Team, channel *model.Channel, user *model.User, removeDeactivatedMembers bool) *model.AppError
	MoveCommand(team *model.Team, command *model.Command) *model.AppError
	MoveFile(oldPath, newPath string) *model.AppError
	NewClusterDiscoveryService() *ClusterDiscoveryService
	NewEmailTemplate(name, locale string) *utils.HTMLTemplate
	NewPluginAPI(manifest *model.Manifest) plugin.API
	NewWebConn(ws *websocket.Conn, session model.Session, t goi18n.TranslateFunc, locale string) *WebConn
	NewWebHub() *Hub
	Notification() einterfaces.NotificationInterface
	NotificationsLog() *mlog.Logger
	OldImportChannel(channel *model.Channel, sChannel SlackChannel, users map[string]*model.User) *model.Channel
	OldImportFile(timestamp time.Time, file io.Reader, teamId string, channelId string, userId string, fileName string) (*model.FileInfo, error)
	OldImportIncomingWebhookPost(post *model.Post, props model.StringInterface) string
	OldImportPost(post *model.Post) string
	OldImportUser(team *model.Team, user *model.User) *model.User
	OpenInteractiveDialog(request model.OpenDialogRequest) *model.AppError
	OriginChecker() func(*http.Request) bool
	OverrideIconURLIfEmoji(post *model.Post)
	PatchBot(botUserId string, botPatch *model.BotPatch) (*model.Bot, *model.AppError)
	PatchChannel(channel *model.Channel, patch *model.ChannelPatch, userId string) (*model.Channel, *model.AppError)
	PatchPost(postId string, patch *model.PostPatch) (*model.Post, *model.AppError)
	PatchRole(role *model.Role, patch *model.RolePatch) (*model.Role, *model.AppError)
	PatchScheme(scheme *model.Scheme, patch *model.SchemePatch) (*model.Scheme, *model.AppError)
	PatchTeam(teamId string, patch *model.TeamPatch) (*model.Team, *model.AppError)
	PatchUser(userId string, patch *model.UserPatch, asAdmin bool) (*model.User, *model.AppError)
	Path() string
	PermanentDeleteAllUsers() *model.AppError
	PermanentDeleteBot(botUserId string) *model.AppError
	PermanentDeleteChannel(channel *model.Channel) *model.AppError
	PermanentDeleteTeam(team *model.Team) *model.AppError
	PermanentDeleteTeamId(teamId string) *model.AppError
	PermanentDeleteUser(user *model.User) *model.AppError
	PluginCommandsForTeam(teamId string) []*model.Command
	PluginContext() *plugin.Context
	PostActionCookieSecret() []byte
	PostAddToChannelMessage(user *model.User, addedUser *model.User, channel *model.Channel, postRootId string) *model.AppError
	PostPatchWithProxyRemovedFromImageURLs(patch *model.PostPatch) *model.PostPatch
	PostUpdateChannelDisplayNameMessage(userId string, channel *model.Channel, oldChannelDisplayName, newChannelDisplayName string) *model.AppError
	PostUpdateChannelHeaderMessage(userId string, channel *model.Channel, oldChannelHeader, newChannelHeader string) *model.AppError
	PostUpdateChannelPurposeMessage(userId string, channel *model.Channel, oldChannelPurpose string, newChannelPurpose string) *model.AppError
	PostWithProxyAddedToImageURLs(post *model.Post) *model.Post
	PostWithProxyRemovedFromImageURLs(post *model.Post) *model.Post
	PreparePostForClient(originalPost *model.Post, isNewPost bool, isEditPost bool) *model.Post
	PreparePostListForClient(originalList *model.PostList) *model.PostList
	ProcessSlackAttachments(attachments []*model.SlackAttachment) []*model.SlackAttachment
	ProcessSlackText(text string) string
	PromoteGuestToUser(user *model.User, requestorId string) *model.AppError
	Publish(message *model.WebSocketEvent)
	PublishSkipClusterSend(message *model.WebSocketEvent)
	PurgeElasticsearchIndexes() *model.AppError
	ReadFile(path string) ([]byte, *model.AppError)
	RecycleDatabaseConnection()
	RegenCommandToken(cmd *model.Command) (*model.Command, *model.AppError)
	RegenOutgoingWebhookToken(hook *model.OutgoingWebhook) (*model.OutgoingWebhook, *model.AppError)
	RegenerateOAuthAppSecret(app *model.OAuthApp) (*model.OAuthApp, *model.AppError)
	RegenerateTeamInviteId(teamId string) (*model.Team, *model.AppError)
	RegisterPluginCommand(pluginId string, command *model.Command) error
	ReloadConfig() error
	RemoveConfigListener(id string)
	RemoveFile(path string) *model.AppError
	RemoveLicense() *model.AppError
	RemoveLicenseListener(id string)
	RemovePlugin(id string) *model.AppError
	RemovePluginFromData(data model.PluginEventData)
	RemoveSamlIdpCertificate() *model.AppError
	RemoveSamlPrivateCertificate() *model.AppError
	RemoveSamlPublicCertificate() *model.AppError
	RemoveTeamIcon(teamId string) *model.AppError
	RemoveTeamMemberFromTeam(teamMember *model.TeamMember, requestorId string) *model.AppError
	RemoveUserFromChannel(userIdToRemove string, removerUserId string, channel *model.Channel) *model.AppError
	RemoveUserFromTeam(teamId string, userId string, requestorId string) *model.AppError
	RenameChannel(channel *model.Channel, newChannelName string, newDisplayName string) (*model.Channel, *model.AppError)
	RenameTeam(team *model.Team, newTeamName string, newDisplayName string) (*model.Team, *model.AppError)
	RequestId() string
	ResetPasswordFromToken(userSuppliedTokenString, newPassword string) *model.AppError
	ResetPermissionsSystem() *model.AppError
	RestoreChannel(channel *model.Channel, userId string) (*model.Channel, *model.AppError)
	RestoreTeam(teamId string) *model.AppError
	RestrictUsersGetByPermissions(userId string, options *model.UserGetOptions) (*model.UserGetOptions, *model.AppError)
	RestrictUsersSearchByPermissions(userId string, options *model.UserSearchOptions) (*model.UserSearchOptions, *model.AppError)
	RevokeAccessToken(token string) *model.AppError
	RevokeAllSessions(userId string) *model.AppError
	RevokeSession(session *model.Session) *model.AppError
	RevokeSessionById(sessionId string) *model.AppError
	RevokeSessionsForDeviceId(userId string, deviceId string, currentSessionId string) *model.AppError
	RevokeSessionsFromAllUsers() *model.AppError
	RevokeUserAccessToken(token *model.UserAccessToken) *model.AppError
	RolesGrantPermission(roleNames []string, permissionId string) bool
	Saml() einterfaces.SamlInterface
	SanitizeProfile(user *model.User, asAdmin bool)
	SanitizeTeam(session model.Session, team *model.Team) *model.Team
	SanitizeTeams(session model.Session, teams []*model.Team) []*model.Team
	SaveAndBroadcastStatus(status *model.Status)
	SaveBrandImage(imageData *multipart.FileHeader) *model.AppError
	SaveComplianceReport(job *model.Compliance) (*model.Compliance, *model.AppError)
	SaveConfig(newCfg *model.Config, sendConfigChangeClusterMessage bool) *model.AppError
	SaveLicense(licenseBytes []byte) (*model.License, *model.AppError)
	SaveReactionForPost(reaction *model.Reaction) (*model.Reaction, *model.AppError)
	SaveUserTermsOfService(userId, termsOfServiceId string, accepted bool) *model.AppError
	SchemesIterator(batchSize int) func() []*model.Scheme
	SearchAllChannels(term string, opts model.ChannelSearchOpts) (*model.ChannelListWithTeamData, int64, *model.AppError)
	SearchAllTeams(searchOpts *model.TeamSearch) ([]*model.Team, int64, *model.AppError)
	SearchArchivedChannels(teamId string, term string, userId string) (*model.ChannelList, *model.AppError)
	SearchChannels(teamId string, term string) (*model.ChannelList, *model.AppError)
	SearchChannelsForUser(userId, teamId, term string) (*model.ChannelList, *model.AppError)
	SearchChannelsUserNotIn(teamId string, userId string, term string) (*model.ChannelList, *model.AppError)
	SearchEmoji(name string, prefixOnly bool, limit int) ([]*model.Emoji, *model.AppError)
	SearchGroupChannels(userId, term string) (*model.ChannelList, *model.AppError)
	SearchPostsInTeam(teamId string, paramsList []*model.SearchParams) (*model.PostList, *model.AppError)
	SearchPostsInTeamForUser(terms string, userId string, teamId string, isOrSearch bool, includeDeletedChannels bool, timeZoneOffset int, page, perPage int) (*model.PostSearchResults, *model.AppError)
	SearchPrivateTeams(term string) ([]*model.Team, *model.AppError)
	SearchPublicTeams(term string) ([]*model.Team, *model.AppError)
	SearchUserAccessTokens(term string) ([]*model.UserAccessToken, *model.AppError)
	SearchUsers(props *model.UserSearch, options *model.UserSearchOptions) ([]*model.User, *model.AppError)
	SearchUsersInChannel(channelId string, term string, options *model.UserSearchOptions) ([]*model.User, *model.AppError)
	SearchUsersInTeam(teamId, term string, options *model.UserSearchOptions) ([]*model.User, *model.AppError)
	SearchUsersNotInChannel(teamId string, channelId string, term string, options *model.UserSearchOptions) ([]*model.User, *model.AppError)
	SearchUsersNotInTeam(notInTeamId string, term string, options *model.UserSearchOptions) ([]*model.User, *model.AppError)
	SearchUsersWithoutTeam(term string, options *model.UserSearchOptions) ([]*model.User, *model.AppError)
	SendAckToPushProxy(ack *model.PushNotificationAck) error
	SendAutoResponse(channel *model.Channel, receiver *model.User) (bool, *model.AppError)
	SendAutoResponseIfNecessary(channel *model.Channel, sender *model.User) (bool, *model.AppError)
	SendChangeUsernameEmail(oldUsername, newUsername, email, locale, siteURL string) *model.AppError
	SendDailyDiagnostics()
	SendDeactivateAccountEmail(email string, locale, siteURL string) *model.AppError
	SendDiagnostic(event string, properties map[string]interface{})
	SendEmailChangeEmail(oldEmail, newEmail, locale, siteURL string) *model.AppError
	SendEmailChangeVerifyEmail(newUserEmail, locale, siteURL, token string) *model.AppError
	SendEmailVerification(user *model.User, newEmail string) *model.AppError
	SendEphemeralPost(userId string, post *model.Post) *model.Post
	SendGuestInviteEmails(team *model.Team, channels []*model.Channel, senderName string, senderUserId string, invites []string, siteURL string, message string)
	SendInviteEmails(team *model.Team, senderName string, senderUserId string, invites []string, siteURL string)
	SendMail(to, subject, htmlBody string) *model.AppError
	SendMailWithEmbeddedFiles(to, subject, htmlBody string, embeddedFiles map[string]io.Reader) *model.AppError
	SendMfaChangeEmail(email string, activated bool, locale, siteURL string) *model.AppError
	SendNotificationMail(to, subject, htmlBody string) *model.AppError
	SendNotifications(post *model.Post, team *model.Team, channel *model.Channel, sender *model.User, parentPostList *model.PostList) ([]string, error)
	SendPasswordChangeEmail(email, method, locale, siteURL string) *model.AppError
	SendPasswordReset(email string, siteURL string) (bool, *model.AppError)
	SendPasswordResetEmail(email string, token *model.Token, locale, siteURL string) (bool, *model.AppError)
	SendSignInChangeEmail(email, method, locale, siteURL string) *model.AppError
	SendUserAccessTokenAddedEmail(email, locale, siteURL string) *model.AppError
	SendVerifyEmail(userEmail, locale, siteURL, token string) *model.AppError
	SendWelcomeEmail(userId string, email string, verified bool, locale, siteURL string) *model.AppError
	ServeInterPluginRequest(w http.ResponseWriter, r *http.Request, sourcePluginId, destinationPluginId string)
	ServePluginPublicRequest(w http.ResponseWriter, r *http.Request)
	ServePluginRequest(w http.ResponseWriter, r *http.Request)
	ServerBusyStateChanged(sbs *model.ServerBusyState)
	Session() *model.Session
	SessionCacheLength() int
	SessionHasPermissionTo(session model.Session, permission *model.Permission) bool
	SessionHasPermissionToChannel(session model.Session, channelId string, permission *model.Permission) bool
	SessionHasPermissionToChannelByPost(session model.Session, postId string, permission *model.Permission) bool
	SessionHasPermissionToManageBot(session model.Session, botUserId string) *model.AppError
	SessionHasPermissionToTeam(session model.Session, teamId string, permission *model.Permission) bool
	SessionHasPermissionToUser(session model.Session, userId string) bool
	SessionHasPermissionToUserOrBot(session model.Session, userId string) bool
	SetAcceptLanguage(s string)
	SetActiveChannel(userId string, channelId string) *model.AppError
	SetAutoResponderStatus(user *model.User, oldNotifyProps model.StringMap)
	SetBotIconImage(botUserId string, file io.ReadSeeker) *model.AppError
	SetBotIconImageFromMultiPartFile(botUserId string, imageData *multipart.FileHeader) *model.AppError
	SetClientLicense(m map[string]string)
	SetContext(c context.Context)
	SetDefaultProfileImage(user *model.User) *model.AppError
	SetDiagnosticId(id string)
	SetIpAddress(s string)
	SetLicense(license *model.License) bool
	SetPath(s string)
	SetPhase2PermissionsMigrationStatus(isComplete bool) error
	SetPluginKey(pluginId string, key string, value []byte) *model.AppError
	SetPluginKeyWithExpiry(pluginId string, key string, value []byte, expireInSeconds int64) *model.AppError
	SetPluginKeyWithOptions(pluginId string, key string, value []byte, options model.PluginKVSetOptions) (bool, *model.AppError)
	SetPluginsEnvironment(pluginsEnvironment *plugin.Environment)
	SetProfileImage(userId string, imageData *multipart.FileHeader) *model.AppError
	SetProfileImageFromFile(userId string, file io.Reader) *model.AppError
	SetProfileImageFromMultiPartFile(userId string, file multipart.File) *model.AppError
	SetRequestId(s string)
	SetSamlIdpCertificateFromMetadata(data []byte) *model.AppError
	SetServer(srv *Server)
	SetSession(s *model.Session)
	SetStatusAwayIfNeeded(userId string, manual bool)
	SetStatusDoNotDisturb(userId string)
	SetStatusLastActivityAt(userId string, activityAt int64)
	SetStatusOffline(userId string, manual bool)
	SetStatusOnline(userId string, manual bool)
	SetStatusOutOfOffice(userId string)
	SetStore(s store.Store)
	SetT(t goi18n.TranslateFunc)
	SetTeamIcon(teamId string, imageData *multipart.FileHeader) *model.AppError
	SetTeamIconFromFile(team *model.Team, file io.Reader) *model.AppError
	SetTeamIconFromMultiPartFile(teamId string, file multipart.File) *model.AppError
	SetUserAgent(s string)
	SetupInviteEmailRateLimiting() error
	ShutDownPlugins()
	Shutdown()
	SlackAddBotUser(teamId string, log *bytes.Buffer) *model.User
	SlackAddChannels(teamId string, slackchannels []SlackChannel, posts map[string][]SlackPost, users map[string]*model.User, uploads map[string]*zip.File, botUser *model.User, importerLog *bytes.Buffer) map[string]*model.Channel
	SlackAddPosts(teamId string, channel *model.Channel, posts []SlackPost, users map[string]*model.User, uploads map[string]*zip.File, botUser *model.User)
	SlackAddUsers(teamId string, slackusers []SlackUser, importerLog *bytes.Buffer) map[string]*model.User
	SlackImport(fileData multipart.File, fileSize int64, teamID string) (*model.AppError, *bytes.Buffer)
	SlackUploadFile(slackPostFile *SlackFile, uploads map[string]*zip.File, teamId string, channelId string, userId string, slackTimestamp string) (*model.FileInfo, bool)
	SoftDeleteTeam(teamId string) *model.AppError
	Srv() *Server
	StartPushNotificationsHubWorkers()
	StopPushNotificationsHubWorkers()
	Store() store.Store
	SubmitInteractiveDialog(request model.SubmitDialogRequest) (*model.SubmitDialogResponse, *model.AppError)
	SwitchEmailToLdap(email, password, code, ldapLoginId, ldapPassword string) (string, *model.AppError)
	SwitchEmailToOAuth(w http.ResponseWriter, r *http.Request, email, password, code, service string) (string, *model.AppError)
	SwitchLdapToEmail(ldapPassword, code, email, newPassword string) (string, *model.AppError)
	SwitchOAuthToEmail(email, password, requesterId string) (string, *model.AppError)
	SyncLdap()
	SyncPlugins() *model.AppError
	SyncPluginsActiveState()
	SyncRolesAndMembership(syncableID string, syncableType model.GroupSyncableType)
	SyncSyncableRoles(syncableID string, syncableType model.GroupSyncableType) *model.AppError
	T(translationID string, args ...interface{}) string
	TeamMembersMinusGroupMembers(teamID string, groupIDs []string, page, perPage int) ([]*model.UserWithGroups, int64, *model.AppError)
	TeamMembersToAdd(since int64, teamID *string) ([]*model.UserTeamIDPair, *model.AppError)
	TeamMembersToRemove(teamID *string) ([]*model.TeamMember, *model.AppError)
	TestElasticsearch(cfg *model.Config) *model.AppError
	TestEmail(userId string, cfg *model.Config) *model.AppError
	TestLdap() *model.AppError
	TestSiteURL(siteURL string) *model.AppError
	Timezones() *timezones.Timezones
	ToggleMuteChannel(channelId string, userId string) *model.ChannelMember
	TotalWebsocketConnections() int
	TriggerWebhook(payload *model.OutgoingWebhookPayload, hook *model.OutgoingWebhook, post *model.Post, channel *model.Channel)
	UnregisterPluginCommand(pluginId, teamId, trigger string)
	UnregisterPluginCommands(pluginId string)
	UpdateActive(user *model.User, active bool) (*model.User, *model.AppError)
	UpdateBotActive(botUserId string, active bool) (*model.Bot, *model.AppError)
	UpdateBotOwner(botUserId, newOwnerId string) (*model.Bot, *model.AppError)
	UpdateChannel(channel *model.Channel) (*model.Channel, *model.AppError)
	UpdateChannelLastViewedAt(channelIds []string, userId string) *model.AppError
	UpdateChannelMemberNotifyProps(data map[string]string, channelId string, userId string) (*model.ChannelMember, *model.AppError)
	UpdateChannelMemberRoles(channelId string, userId string, newRoles string) (*model.ChannelMember, *model.AppError)
	UpdateChannelMemberSchemeRoles(channelId string, userId string, isSchemeGuest bool, isSchemeUser bool, isSchemeAdmin bool) (*model.ChannelMember, *model.AppError)
	UpdateChannelPrivacy(oldChannel *model.Channel, user *model.User) (*model.Channel, *model.AppError)
	UpdateChannelScheme(channel *model.Channel) (*model.Channel, *model.AppError)
	UpdateCommand(oldCmd, updatedCmd *model.Command) (*model.Command, *model.AppError)
	UpdateConfig(f func(*model.Config))
	UpdateEphemeralPost(userId string, post *model.Post) *model.Post
	UpdateGroup(group *model.Group) (*model.Group, *model.AppError)
	UpdateGroupSyncable(groupSyncable *model.GroupSyncable) (*model.GroupSyncable, *model.AppError)
	UpdateIncomingWebhook(oldHook, updatedHook *model.IncomingWebhook) (*model.IncomingWebhook, *model.AppError)
	UpdateLastActivityAtIfNeeded(session model.Session)
	UpdateMfa(activate bool, userId, token string) *model.AppError
	UpdateMobileAppBadge(userId string)
	UpdateMobileAppBadgeSync(userId string) *model.AppError
	UpdateOAuthUserAttrs(userData io.Reader, user *model.User, provider einterfaces.OauthProvider, service string) *model.AppError
	UpdateOauthApp(oldApp, updatedApp *model.OAuthApp) (*model.OAuthApp, *model.AppError)
	UpdateOutgoingWebhook(oldHook, updatedHook *model.OutgoingWebhook) (*model.OutgoingWebhook, *model.AppError)
	UpdatePassword(user *model.User, newPassword string) *model.AppError
	UpdatePasswordAsUser(userId, currentPassword, newPassword string) *model.AppError
	UpdatePasswordByUserIdSendEmail(userId, newPassword, method string) *model.AppError
	UpdatePasswordSendEmail(user *model.User, newPassword, method string) *model.AppError
	UpdatePost(post *model.Post, safeUpdate bool) (*model.Post, *model.AppError)
	UpdatePreferences(userId string, preferences model.Preferences) *model.AppError
	UpdateRole(role *model.Role) (*model.Role, *model.AppError)
	UpdateScheme(scheme *model.Scheme) (*model.Scheme, *model.AppError)
	UpdateSessionsIsGuest(userId string, isGuest bool)
	UpdateTeam(team *model.Team) (*model.Team, *model.AppError)
	UpdateTeamMemberRoles(teamId string, userId string, newRoles string) (*model.TeamMember, *model.AppError)
	UpdateTeamMemberSchemeRoles(teamId string, userId string, isSchemeGuest bool, isSchemeUser bool, isSchemeAdmin bool) (*model.TeamMember, *model.AppError)
	UpdateTeamPrivacy(teamId string, teamType string, allowOpenInvite bool) *model.AppError
	UpdateTeamScheme(team *model.Team) (*model.Team, *model.AppError)
	UpdateUser(user *model.User, sendNotifications bool) (*model.User, *model.AppError)
	UpdateUserActive(userId string, active bool) *model.AppError
	UpdateUserAsUser(user *model.User, asAdmin bool) (*model.User, *model.AppError)
	UpdateUserAuth(userId string, userAuth *model.UserAuth) (*model.UserAuth, *model.AppError)
	UpdateUserNotifyProps(userId string, props map[string]string) (*model.User, *model.AppError)
	UpdateUserRoles(userId string, newRoles string, sendWebSocketEvent bool) (*model.User, *model.AppError)
	UpdateWebConnUserActivity(session model.Session, activityAt int64)
	UploadEmojiImage(id string, imageData *multipart.FileHeader) *model.AppError
	UploadFile(data []byte, channelId string, filename string) (*model.FileInfo, *model.AppError)
	UploadFileX(channelId, name string, input io.Reader, opts ...func(*UploadFileTask)) (*model.FileInfo, *model.AppError)
	UploadFiles(teamId string, channelId string, userId string, files []io.ReadCloser, filenames []string, clientIds []string, now time.Time) (*model.FileUploadResponse, *model.AppError)
	UploadMultipartFiles(teamId string, channelId string, userId string, fileHeaders []*multipart.FileHeader, clientIds []string, now time.Time) (*model.FileUploadResponse, *model.AppError)
	UpsertGroupMember(groupID string, userID string) (*model.GroupMember, *model.AppError)
	UpsertGroupSyncable(groupSyncable *model.GroupSyncable) (*model.GroupSyncable, *model.AppError)
	UserAgent() string
	UserCanSeeOtherUser(userId string, otherUserId string) (bool, *model.AppError)
	UserIsInAdminRoleGroup(userID, syncableID string, syncableType model.GroupSyncableType) (bool, *model.AppError)
	ValidateAndSetLicenseBytes(b []byte)
	VerifyEmailFromToken(userSuppliedTokenString string) *model.AppError
	VerifyPlugin(plugin, signature io.ReadSeeker) *model.AppError
	VerifyUserEmail(userId, email string) *model.AppError
	ViewChannel(view *model.ChannelView, userId string, currentSessionId string) (map[string]int64, *model.AppError)
	WaitForChannelMembership(channelId string, userId string)
	WriteFile(fr io.Reader, path string) (int64, *model.AppError)
}
