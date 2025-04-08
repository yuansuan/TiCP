package guacamole

type ConnectArgs struct {
	ConnectArgsInToken
	ScreenWidth  int `form:"screen_width" json:"screen_width"`
	ScreenHeight int `form:"screen_height" json:"screen_height"`
}

type ConnectArgsInToken struct {
	GuacadAddr         string `form:"guacad_addr" json:"guacad_addr"`
	AssetProtocol      string `form:"asset_protocol" json:"asset_protocol"` // rdp
	AssetHost          string `form:"asset_host" json:"asset_host"`
	AssetPort          string `form:"asset_port" json:"asset_port"`
	AssetUser          string `form:"asset_user" json:"asset_user"`
	AssetPassword      string `form:"asset_password" json:"asset_password"`
	AssetSecurity      string `form:"asset_security" json:"asset_security"`
	AssetRemoteApp     string `form:"asset_remote_app" json:"asset_remote_app"`
	AssetRemoteAppArgs string `form:"asset_remote_app_args" json:"asset_remote_app_args"`
	AssetRemoteAppDir  string `form:"asset_remote_app_dir" json:"asset_remote_app_dir"`
	StorageID          string `form:"storage_id" json:"storage_id"`
	DisableGfx         string `form:"disable_gfx" json:"disable_gfx"`
	ResizeMethod       string `form:"resize_method" json:"resize_method"` // default: "display-update", another choice is "reconnect"
}
