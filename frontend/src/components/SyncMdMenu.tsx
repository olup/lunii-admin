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
  MdMenu,
  MdOutlineCloud,
  MdOutlineComputer,
} from "react-icons/md";
import { ListPacks, SyncLuniiStoreMetadata } from "../../wailsjs/go/main/App";

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

  return (
    <Menu>
      <MenuButton
        as={Button}
        rightIcon={<MdMenu />}
        variant="outline"
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
        <MenuItem icon={<MdOutlineComputer />}>From default studio DB</MenuItem>
        <MenuItem icon={<MdOutlineComputer />}>From custom studio db</MenuItem>
      </MenuList>
    </Menu>
  );
};
