// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {uuid} from '../models';
import {main} from '../models';
import {lunii} from '../models';

export function ChangePackOrder(arg1:uuid.UUID,arg2:number):Promise<string>;

export function CheckForUpdate():Promise<main.CheckUpdateResponse>;

export function CreatePack(arg1:string,arg2:string):Promise<string>;

export function GetDeviceInfos():Promise<lunii.Device>;

export function GetInfos():Promise<main.Infos>;

export function InstallPack(arg1:string):Promise<string>;

export function ListPacks():Promise<Array<lunii.Metadata>>;

export function OpenDirectory(arg1:string):Promise<string>;

export function OpenFile(arg1:string):Promise<string>;

export function OpenFiles(arg1:string):Promise<Array<string>>;

export function RemovePack(arg1:uuid.UUID):Promise<boolean>;

export function SaveFile(arg1:string,arg2:string,arg3:string):Promise<string>;

export function SyncLuniiStoreMetadata(arg1:Array<uuid.UUID>):Promise<string>;

export function SyncStudioMetadata(arg1:Array<uuid.UUID>,arg2:string):Promise<string>;
