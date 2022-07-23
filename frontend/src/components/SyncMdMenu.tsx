import {
  Menu,
  MenuButton,
  Button,
  MenuList,
  MenuItem,
  useToast,
} from "@chakra-ui/react";
import { useQuery } from "@tanstack/react-query";
import {
  MdArrowDownward,
  MdCloud,
  MdDownload,
  MdMenu,
  MdOutlineCloud,
  MdOutlineComputer,
  MdSettings,
} from "react-icons/md";
import {
  ListPacks,
  OpenDirectory,
  OpenFile,
  SyncLuniiStoreMetadata,
  SyncStudioMetadata,
} from "../../wailsjs/go/main/App";

export const SyncMdMenu = () => {
  const { data: packs, refetch } = useQuery(["packs"], ListPacks);
  const toast = useToast();

  const handleSyncLuniiStoreMetadata = async () => {
    const uuids = packs?.map((p) => p.uuid) || [];
    await SyncLuniiStoreMetadata(uuids);
    toast({
      title: "Metadata have been downloaded",
      status: "success",
    });
    await refetch();
  };

  const handleSyncStudioMetadata = async (customPath = false) => {
    const uuids = packs?.map((p) => p.uuid) || [];
    let dbPath = "";

    if (customPath) {
      dbPath = await OpenFile("Select Studio DB location");
    }

    await SyncStudioMetadata(uuids, dbPath);
    toast({
      title: "Metadata have been downloaded",
      status: "success",
    });
    await refetch();
  };

  return (
    <Menu>
      <MenuButton
        as={Button}
        leftIcon={<MdDownload />}
        variant="ghost"
        colorScheme="linkedin"
      >
        Sync Metadata
      </MenuButton>
      <MenuList>
        <MenuItem
          icon={<MdOutlineCloud />}
          onClick={handleSyncLuniiStoreMetadata}
        >
          From lunii store
        </MenuItem>
        <MenuItem
          icon={<MdOutlineComputer />}
          onClick={() => handleSyncStudioMetadata()}
        >
          From default studio DB
        </MenuItem>
        <MenuItem
          icon={<MdOutlineComputer />}
          onClick={() => handleSyncStudioMetadata(true)}
        >
          From custom studio db
        </MenuItem>
      </MenuList>
    </Menu>
  );
};
