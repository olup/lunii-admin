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
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tooltip,
  Tr,
  useToast,
} from "@chakra-ui/react";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import {
  FiArrowDown,
  FiArrowUp,
  FiRefreshCcw,
  FiSmartphone,
  FiUpload,
} from "react-icons/fi";
import { BiPlug } from "react-icons/bi";
import {
  ChangePackOrder,
  GetDeviceInfos,
  InstallPack,
  ListPacks,
} from "../wailsjs/go/main/App";
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
                rightIcon={<FiRefreshCcw />}
                onClick={() => refetchDevice()}
              >
                Refresh
              </Button>
            </Tooltip>
            <NewPackModal />
          </Box>
          <Center opacity={0.3} fontSize={20} h={500} flexDirection="column">
            <Icon as={BiPlug}></Icon>

            <Text>No device connected</Text>
          </Center>
        </>
      )}

      {device && (
        <Box pt="56px">
          <Box
            display="flex"
            position="fixed"
            backgroundColor="white"
            zIndex={2}
            top={0}
            right={0}
            left={0}
            p={2}
          >
            <Box>
              <Popover placement="bottom-start" closeOnBlur={false}>
                <PopoverTrigger>
                  <Button colorScheme="linkedin" leftIcon={<FiSmartphone />}>
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
            <Box ml={2}>
              <SyncMdMenu />
            </Box>
            <Tooltip label="Install a STUdio story pack to your device">
              <Button
                variant="ghost"
                colorScheme="linkedin"
                leftIcon={<FiUpload />}
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
                      icon={<FiArrowUp />}
                      mr={1}
                      onClick={() => handleChangePackOrder(p.uuid, i - 1)}
                    />
                    <IconButton
                      size="xs"
                      aria-label="down"
                      icon={<FiArrowDown />}
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
                    <DetailsModal uuid={p.uuid} />
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
