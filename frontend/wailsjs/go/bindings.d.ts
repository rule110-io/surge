interface go {
  "surge": {
    "MiddlewareFunctions": {
		DownloadFile(arg1:string):Promise<boolean>
		GetFileChunkMap(arg1:string,arg2:number):Promise<string>
		GetFileDetails(arg1:string):Promise<FileDetails>
		GetLocalFiles(arg1:string,arg2:FileFilterState,arg3:string,arg4:boolean,arg5:number,arg6:number):Promise<PagedQueryResult>
		GetOfficialTopicName():Promise<string>
		GetPublicKey():Promise<string>
		GetRemoteFiles(arg1:string,arg2:string,arg3:string,arg4:boolean,arg5:number,arg6:number):Promise<PagedQueryRemoteResult>
		GetTopicDetails(arg1:string):Promise<TopicInfo>
		GetTopicSubscriptions():Promise<Array<string>>
		OpenFile(arg1:string):Promise<void>
		OpenFolder(arg1:string):Promise<void>
		OpenLink(arg1:string):Promise<void>
		OpenLog():Promise<void>
		ReadSetting(arg1:string):Promise<string>
		RemoveFile(arg1:string,arg2:boolean):Promise<boolean>
		SeedFile(arg1:string):Promise<boolean>
		SetDownloadPause(arg1:Array<string>,arg2:boolean):Promise<void>
		StartDownloadMagnetLinks(arg1:string):Promise<boolean>
		SubscribeToTopic(arg1:string):Promise<void>
		UnsubscribeFromTopic(arg1:string):Promise<void>
		WriteSetting(arg1:string,arg2:string):Promise<boolean>
    },
  }

}

declare global {
	interface Window {
		go: go;
	}
}
