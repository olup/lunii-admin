import {
  Box,
  Button,
  Center,
  Code,
  Flex,
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
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import {
  FiArrowDown,
  FiArrowUp,
  FiMenu,
  FiRefreshCcw,
  FiSmartphone,
  FiUpload,
} from "react-icons/fi";
import { BiMenu, BiPlug } from "react-icons/bi";
import { MdDragIndicator } from "react-icons/md";
import {
  ChangePackOrder,
  GetDeviceInfos,
  InstallPack,
  ListPacks,
  OpenFile,
} from "../wailsjs/go/main/App";
import { DetailsModal } from "./components/DetailsModal";
import { IsInstallingModal } from "./components/IsInstallingModal";
import { NewPackModal } from "./components/NewPackModal";
import { PackTag } from "./components/PackTag";
import { SyncMdMenu } from "./components/SyncMdMenu";
import { DragDropContext, Draggable, Droppable } from "react-beautiful-dnd";
import { lunii } from "../wailsjs/go/models";
import { useChangePackOrder } from "./hooks/useChangePackOrder";

function App() {
  const {
    data: device,
    refetch: refetchDevice,
    isLoading: isLoadingDevice,
  } = useQuery(["device"], GetDeviceInfos);
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
    const packPath = await OpenFile("Select your pack");
    if (!packPath) return;

    setIsInstalling(true);
    try {
      await InstallPack(packPath);
      toast({
        title: "The pack was installed on the device",
        status: "success",
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

  const { mutate: handleChangePackOrder } = useChangePackOrder();

  if (isLoadingDevice) return null;

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
            <Button
              variant="ghost"
              colorScheme="linkedin"
              leftIcon={<FiUpload />}
              onClick={handleInstallStory}
              ml={2}
            >
              Install pack
            </Button>
            <NewPackModal />
          </Box>

          <Box>
            <Flex py={3}>
              <Box width={50}></Box>
              <Box flex={1}>Title</Box>
              <Box width={100}></Box>
              <Box width={50}></Box>
            </Flex>
            <DragDropContext
              onDragEnd={(a) => {
                if (!a.destination) return;
                handleChangePackOrder({
                  id: a.draggableId,
                  position: a.destination?.index,
                });
              }}
            >
              <Droppable droppableId="droppable">
                {(provided) => (
                  <Box {...provided.droppableProps} ref={provided.innerRef}>
                    {packs?.map((p, i) => (
                      <Draggable
                        key={p.uuid as any}
                        draggableId={p.uuid as any}
                        index={i}
                      >
                        {(provided) => (
                          <Flex
                            ref={provided.innerRef}
                            {...provided.draggableProps}
                            backgroundColor="white"
                            borderBottom="1px solid #eee"
                            borderTop="1px solid #eee"
                            mt="-1px"
                            py={2}
                            alignItems="center"
                          >
                            <Box width={50}>
                              <IconButton
                                size="xs"
                                aria-label="up"
                                icon={<MdDragIndicator />}
                                variant="ghost"
                                mr={1}
                                {...provided.dragHandleProps}
                              />
                            </Box>

                            <Box
                              flex={1}
                              fontWeight={p.title && "bold"}
                              opacity={p.title ? 1 : 0.5}
                            >
                              {p.title || p.uuid}
                            </Box>
                            <Box width={100}>
                              <PackTag metadata={p} />
                            </Box>
                            <Box width={50}>
                              <DetailsModal uuid={p.uuid} />
                            </Box>
                          </Flex>
                        )}
                      </Draggable>
                    ))}
                  </Box>
                )}
              </Droppable>
            </DragDropContext>
          </Box>
        </Box>
      )}
      <IsInstallingModal isOpen={isInstalling} />
    </Box>
  );
}

export default App;
