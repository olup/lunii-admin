// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {uuid} from '../models';
import {lunii} from '../models';

export function RemovePack(arg1:uuid.UUID):Promise<boolean|Error>;

export function SyncLuniiStoreMetadata(arg1:Array<uuid.UUID>):Promise<string|Error>;

export function SyncStudioMetadata(arg1:Array<uuid.UUID>,arg2:string):Promise<string|Error>;

export function ChangePackOrder(arg1:uuid.UUID,arg2:number):Promise<string|Error>;

export function InstallPack():Promise<string|Error>;

export function ListPacks():Promise<Array<lunii.Metadata>>;

export function OpenDirectory(arg1:string):Promise<string>;

export function CreatePack(arg1:string,arg2:string):Promise<string|Error>;

export function GetDeviceInfos():Promise<lunii.Device>;

export function OpenFile(arg1:string):Promise<string>;

export function SaveFile(arg1:string,arg2:string,arg3:string):Promise<string>;
