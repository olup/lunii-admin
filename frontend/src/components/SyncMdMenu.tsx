import {
  Menu,
  MenuButton,
  Button,
  MenuList,
  MenuItem,
  useToast,
} from "@chakra-ui/react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { FiCloud, FiDownload, FiHardDrive } from "react-icons/fi";
import {
  ListPacks,
  OpenDirectory,
  OpenFile,
  SyncLuniiStoreMetadata,
  SyncStudioMetadata,
} from "../../wailsjs/go/main/App";

export const SyncMdMenu = () => {
  const { data: packs } = useQuery(["packs"], ListPacks);
  const queryClient = useQueryClient();
  const toast = useToast();

  const handleSyncLuniiStoreMetadata = async () => {
    try {
      const uuids = packs?.map((p) => p.uuid) || [];
      await SyncLuniiStoreMetadata(uuids);
      toast({
        title: "Metadata have been synced",
        status: "success",
      });
    } catch (err) {
      toast({
        title: "Could not sync metadata",
        status: "error",
        description: err as string,
      });
    }

    await queryClient.invalidateQueries(["packs"]);
  };

  const handleSyncStudioMetadata = async (customPath = false) => {
    try {
      const uuids = packs?.map((p) => p.uuid) || [];
      let dbPath = "";

      if (customPath) {
        dbPath = await OpenFile("Select Studio DB location");
      }

      await SyncStudioMetadata(uuids, dbPath);
      toast({
        title: "Metadata have been synced",
        status: "success",
      });
    } catch (err) {
      toast({
        title: "Could not sync metadata",
        status: "error",
        description: err as string,
      });
    }
    await queryClient.invalidateQueries(["packs"]);
  };

  return (
    <Menu>
      <MenuButton
        as={Button}
        leftIcon={<FiDownload />}
        variant="ghost"
        colorScheme="linkedin"
      >
        Sync Metadata
      </MenuButton>
      <MenuList>
        <MenuItem icon={<FiCloud />} onClick={handleSyncLuniiStoreMetadata}>
          From lunii store
        </MenuItem>
        <MenuItem
          icon={<FiHardDrive />}
          onClick={() => handleSyncStudioMetadata()}
        >
          From default studio DB
        </MenuItem>
        <MenuItem
          icon={<FiHardDrive />}
          onClick={() => handleSyncStudioMetadata(true)}
        >
          From custom studio db
        </MenuItem>
      </MenuList>
    </Menu>
  );
};
