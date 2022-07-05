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
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import {
  MdCreate,
  MdDelete,
  MdDevices,
  MdRefresh,
  MdUpload,
} from "react-icons/md";
import {
  GetDeviceInfos,
  ListPacks,
  InstallPack,
  RemovePack,
} from "../wailsjs/go/main/App";
import { lunii } from "../wailsjs/go/models";
import { NewStoryModal } from "./components/NewStoryModal";

function App() {
  const [device, setDevice] = useState<lunii.Device>();
  const [loading, setLoading] = useState(false);
  const [packs, setPacks] = useState<lunii.Metadata[]>([]);

  const loadDevice = () => {
    setLoading(true);
    GetDeviceInfos()
      .then(setDevice)
      .finally(() => setLoading(false));
  };

  const loadPacks = () => {
    ListPacks().then(setPacks);
  };

  useEffect(() => {
    loadDevice();
  }, []);

  useEffect(() => {
    loadPacks();
  }, [device]);

  const handleInstallStory = async () => {
    await InstallPack();
    await loadPacks();
  };

  return (
    <Box id="App" p={3}>
      {!device && (
        <>
          <Box display="flex">
            <Button
              colorScheme="orange"
              rightIcon={<MdRefresh />}
              onClick={loadDevice}
              isLoading={loading}
            >
              Refresh
            </Button>
            <NewStoryModal />
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
                    Searial Number <Code>{device.serialNumber}</Code>
                    Version{" "}
                    <Code>
                      {device.firmwareVersionMajor}.
                      {device.firmwareVersionMinor}
                    </Code>
                  </Box>
                </PopoverBody>
              </PopoverContent>
            </Popover>
            <Tooltip label="Install a STUdio story pack to your device">
              <Button
                variant="outline"
                colorScheme="linkedin"
                rightIcon={<MdUpload />}
                onClick={handleInstallStory}
                ml={2}
              >
                Install story
              </Button>
            </Tooltip>
            <Tooltip label="Create a STUdio story pack from a directory">
              <Button
                variant="outline"
                colorScheme="linkedin"
                rightIcon={<MdCreate />}
                ml={2}
              >
                Create story
              </Button>
            </Tooltip>
          </Box>

          <Box>
            <Table>
              <Thead>
                <Tr>
                  <Th fontWeight="bold">Title</Th>
                  <Th>Ref</Th>
                  <Th>Uuid</Th>
                  <Th></Th>
                </Tr>
              </Thead>
              {packs.map((p) => (
                <Tbody>
                  <Td fontWeight="bold">{p.title}</Td>
                  <Td>{p.ref}</Td>
                  <Td>{p.uuid}</Td>
                  <Td>
                    <IconButton
                      variant="ghost"
                      aria-label="delete"
                      icon={<MdDelete />}
                      onClick={() => RemovePack(p.ref)}
                    />
                  </Td>
                </Tbody>
              ))}
            </Table>
          </Box>
        </Box>
      )}
    </Box>
  );
}

export default App;
