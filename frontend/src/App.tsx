import {
  Box,
  Button,
  Center,
  Code,
  Icon,
  IconButton,
  Popover,
  PopoverBody,
  PopoverCloseButton,
  PopoverContent,
  PopoverHeader,
  PopoverTrigger,
  Table,
  Tag,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tooltip,
  Tr,
  useToast,
} from "@chakra-ui/react";
import { QueryClient, useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import {
  MdArrowDownward,
  MdArrowUpward,
  MdCreate,
  MdDelete,
  MdDevices,
  MdRefresh,
  MdSync,
  MdUpload,
  MdWarning,
} from "react-icons/md";
import {
  GetDeviceInfos,
  ListPacks,
  InstallPack,
  RemovePack,
  ChangePackOrder,
  SyncLuniiStoreMetadata,
} from "../wailsjs/go/main/App";
import { lunii } from "../wailsjs/go/models";
import { DeleteModal } from "./components/DeleteModal";
import { DetailsModal } from "./components/DetailsModal";
import { IsInstallingModal } from "./components/IsInstallingModal";
import { NewPackModal } from "./components/NewPackModal";
import { PackTag } from "./components/PackTag";
import { SyncMdMenu } from "./components/SyncMdMenu";

function App() {
  const { data: device, refetch: refetchDevice } = useQuery(
    ["device"],
    GetDeviceInfos
  );
  const { data: packs, refetch: refetchPacks } = useQuery(
    ["packs"],
    ListPacks,
    {
      refetchOnMount: true,
      enabled: !!device,
    }
  );

  const [isInstalling, setIsInstalling] = useState(false);
  const toast = useToast();

  const handleInstallStory = async () => {
    setIsInstalling(true);
    try {
      await InstallPack();
      toast({
        title: "The pack was installed on the device",
        status: "success",
        duration: 9000,
        isClosable: true,
      });
      await refetchPacks();
    } catch (e) {
      toast({
        title: "Could not install pack on device",
        status: "warning",
        description:
          "Something went wrong while installing the pack on your device",
      });
    }
    setIsInstalling(false);
  };

  const handleRemovePack = async (uuid: number[]) => {
    await RemovePack(uuid);
    toast({
      title: "The pack was deleted from the device",
      status: "success",
      duration: 9000,
      isClosable: true,
    });
    await refetchPacks();
  };

  const handleChangePackOrder = async (uuid: number[], position: number) => {
    await ChangePackOrder(uuid, position);
    await refetchPacks();
  };

  return (
    <Box id="App" p={3}>
      {!device && (
        <>
          <Box display="flex">
            <Tooltip label="Try refreshing after the device has been mounted by the system">
              <Button
                colorScheme="orange"
                rightIcon={<MdRefresh />}
                onClick={() => refetchDevice()}
              >
                Refresh
              </Button>
            </Tooltip>
            <NewPackModal />
          </Box>
          <Center opacity={0.3} fontSize={30} h={500} flexDirection="column">
            <Icon fontSize={100} as={MdDevices}></Icon>

            <Text>No device connected</Text>
          </Center>
        </>
      )}

      {device && (
        <Box>
          <Box display="flex">
            <Box mr={2}>
              <Popover placement="bottom-start" closeOnBlur={false}>
                <PopoverTrigger>
                  <Button colorScheme="linkedin" leftIcon={<MdDevices />}>
                    My Lunii
                  </Button>
                </PopoverTrigger>
                <PopoverContent>
                  <PopoverCloseButton />
                  <PopoverHeader>Details</PopoverHeader>
                  <PopoverBody>
                    <Box>
                      Serial Number <Code>{device.serialNumber}</Code>
                      Version{" "}
                      <Code>
                        {device.firmwareVersionMajor}.
                        {device.firmwareVersionMinor}
                      </Code>
                    </Box>
                  </PopoverBody>
                </PopoverContent>
              </Popover>
            </Box>
            <SyncMdMenu />
            <Tooltip label="Install a STUdio story pack to your device">
              <Button
                variant="outline"
                colorScheme="linkedin"
                rightIcon={<MdUpload />}
                onClick={handleInstallStory}
                ml={2}
              >
                Install pack
              </Button>
            </Tooltip>
            <NewPackModal />
          </Box>

          <Box>
            <Table>
              <Thead>
                <Tr>
                  <Th></Th>
                  <Th>Title</Th>
                  <Th></Th>
                  <Th></Th>
                </Tr>
              </Thead>
              {packs?.map((p, i) => (
                <Tbody>
                  <Td>
                    <IconButton
                      size="xs"
                      aria-label="up"
                      icon={<MdArrowUpward />}
                      mr={1}
                      onClick={() => handleChangePackOrder(p.uuid, i - 1)}
                    />
                    <IconButton
                      size="xs"
                      aria-label="down"
                      icon={<MdArrowDownward />}
                      onClick={() => handleChangePackOrder(p.uuid, i + 1)}
                    />
                  </Td>

                  <Td
                    fontWeight={p.title && "bold"}
                    opacity={p.title ? 1 : 0.5}
                  >
                    {p.title || p.uuid}
                  </Td>
                  <Td>
                    <PackTag metadata={p} />
                  </Td>
                  <Td>
                    <DetailsModal
                      metadata={p}
                      onDelete={() => handleRemovePack(p.uuid)}
                    />
                  </Td>
                </Tbody>
              ))}
            </Table>
          </Box>
        </Box>
      )}
      <IsInstallingModal isOpen={isInstalling} />
    </Box>
  );
}

export default App;
