import { useState, useEffect } from "react";
import axios from "axios";
import {
  createTheme,
  ThemeProvider,
  CssBaseline,
  Container,
  AppBar,
  Toolbar,
  Typography,
  Tabs,
  Tab,
  Box,
  Button,
  TextField,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Snackbar,
  Alert,
  Card,
  CardContent,
  Grid,
  Chip,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Divider,
} from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";
import AddIcon from "@mui/icons-material/Add";
import SearchIcon from "@mui/icons-material/Search";
import BarChartIcon from "@mui/icons-material/BarChart";
import LibraryMusicIcon from "@mui/icons-material/LibraryMusic";
import QueueMusicIcon from "@mui/icons-material/QueueMusic";
import PlaylistAddIcon from "@mui/icons-material/PlaylistAdd";
import FilterListIcon from "@mui/icons-material/FilterList";

const darkTheme = createTheme({
  palette: {
    mode: "dark",
    primary: { main: "#7c4dff" },
    secondary: { main: "#ff4081" },
    background: { default: "#0a0a0a", paper: "#141414" },
  },
  typography: { fontFamily: "'Inter', sans-serif" },
});

const API = "http://localhost:8080/api";

export default function App() {
  const [tab, setTab] = useState(0);
  const [tracks, setTracks] = useState([]);
  const [filteredTracks, setFilteredTracks] = useState([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState([]);
  const [stats, setStats] = useState(null);
  const [playlists, setPlaylists] = useState([]);
  const [selectedPlaylist, setSelectedPlaylist] = useState(null);
  const [playlistTracks, setPlaylistTracks] = useState([]);
  const [openAdd, setOpenAdd] = useState(false);
  const [openEdit, setOpenEdit] = useState(false);
  const [openAddPlaylist, setOpenAddPlaylist] = useState(false);
  const [openRenamePlaylist, setOpenRenamePlaylist] = useState(false);
  const [openAddToPlaylist, setOpenAddToPlaylist] = useState(false);
  const [selectedTrackForPlaylist, setSelectedTrackForPlaylist] = useState(null);
  const [newPlaylistName, setNewPlaylistName] = useState("");
  const [renameValue, setRenameValue] = useState("");
  const [snackbar, setSnackbar] = useState({ open: false, message: "", severity: "success" });
  const [newTrack, setNewTrack] = useState({ artist: "", track: "", url: "" });
  const [editTrack, setEditTrack] = useState({ track_id: null, artist: "", track: "", url: "" });
  const [filterFrom, setFilterFrom] = useState("");
  const [filterTo, setFilterTo] = useState("");
  const [filterActive, setFilterActive] = useState(false);

  useEffect(() => { fetchTracks(); }, []);

  const fetchTracks = async () => {
    try {
      const res = await axios.get(`${API}/tracks`);
      setTracks(res.data || []);
      setFilteredTracks(res.data || []);
    } catch { showSnackbar("Ошибка загрузки треков", "error"); }
  };

  const fetchStats = async () => {
    try {
      const res = await axios.get(`${API}/statistics`);
      setStats(res.data);
    } catch { showSnackbar("Ошибка загрузки статистики", "error"); }
  };

  const fetchPlaylists = async () => {
    try {
      const res = await axios.get(`${API}/playlists`);
      setPlaylists(res.data || []);
    } catch { showSnackbar("Ошибка загрузки плейлистов", "error"); }
  };

  const fetchPlaylistTracks = async (id) => {
    try {
      const res = await axios.get(`${API}/playlists/${id}/tracks`);
      setPlaylistTracks(res.data || []);
    } catch { showSnackbar("Ошибка загрузки треков плейлиста", "error"); }
  };

  const handleAddTrack = async () => {
    if (!newTrack.artist || !newTrack.track) {
      showSnackbar("Заполните исполнителя и название", "warning");
      return;
    }
    try {
      await axios.post(`${API}/tracks`, newTrack);
      showSnackbar("Трек добавлен!", "success");
      setOpenAdd(false);
      setNewTrack({ artist: "", track: "", url: "" });
      fetchTracks();
    } catch { showSnackbar("Ошибка добавления", "error"); }
  };

  const handleDeleteTrack = async (id) => {
    try {
      await axios.delete(`${API}/tracks/${id}`);
      showSnackbar("Трек удалён", "success");
      fetchTracks();
    } catch { showSnackbar("Ошибка удаления", "error"); }
  };

  const handleEditTrack = async () => {
    try {
      await axios.put(`${API}/tracks/${editTrack.track_id}`, editTrack);
      showSnackbar("Трек обновлён!", "success");
      setOpenEdit(false);
      fetchTracks();
    } catch { showSnackbar("Ошибка обновления", "error"); }
  };

  const handleSearch = async () => {
    if (!searchQuery) { fetchTracks(); return; }
    try {
      const res = await axios.get(`${API}/tracks/search?query=${searchQuery}`);
      setSearchResults(res.data || []);
    } catch { showSnackbar("Ошибка поиска", "error"); }
  };

  const handleFilter = async () => {
    try {
      const params = new URLSearchParams();
      if (filterFrom) params.append("from", filterFrom);
      if (filterTo) params.append("to", filterTo);
      const res = await axios.get(`${API}/tracks/filter?${params.toString()}`);
      setFilteredTracks(res.data || []);
      setFilterActive(true);
      showSnackbar(`Найдено ${res.data?.length || 0} треков`, "success");
    } catch { showSnackbar("Ошибка фильтрации", "error"); }
  };

  const handleClearFilter = () => {
    setFilterFrom("");
    setFilterTo("");
    setFilteredTracks(tracks);
    setFilterActive(false);
  };

  const handleCreatePlaylist = async () => {
    if (!newPlaylistName) { showSnackbar("Введите название", "warning"); return; }
    try {
      await axios.post(`${API}/playlists`, { name: newPlaylistName });
      showSnackbar("Плейлист создан!", "success");
      setOpenAddPlaylist(false);
      setNewPlaylistName("");
      fetchPlaylists();
    } catch { showSnackbar("Ошибка создания плейлиста", "error"); }
  };

  const handleDeletePlaylist = async (id) => {
    try {
      await axios.delete(`${API}/playlists/${id}`);
      showSnackbar("Плейлист удалён", "success");
      if (selectedPlaylist?.playlist_id === id) {
        setSelectedPlaylist(null);
        setPlaylistTracks([]);
      }
      fetchPlaylists();
    } catch { showSnackbar("Ошибка удаления плейлиста", "error"); }
  };

  const handleRenamePlaylist = async () => {
    try {
      await axios.put(`${API}/playlists/${selectedPlaylist.playlist_id}`, { name: renameValue });
      showSnackbar("Плейлист переименован", "success");
      setOpenRenamePlaylist(false);
      fetchPlaylists();
      setSelectedPlaylist({ ...selectedPlaylist, name: renameValue });
    } catch { showSnackbar("Ошибка переименования", "error"); }
  };

  const handleAddTrackToPlaylist = async (playlistId) => {
    try {
      await axios.post(`${API}/playlists/${playlistId}/tracks/${selectedTrackForPlaylist.track_id}`);
      showSnackbar("Трек добавлен в плейлист!", "success");
      setOpenAddToPlaylist(false);
      if (selectedPlaylist?.playlist_id === playlistId) {
        fetchPlaylistTracks(playlistId);
      }
    } catch { showSnackbar("Ошибка добавления в плейлист", "error"); }
  };

  const handleRemoveFromPlaylist = async (trackId) => {
    try {
      await axios.delete(`${API}/playlists/${selectedPlaylist.playlist_id}/tracks/${trackId}`);
      showSnackbar("Трек удалён из плейлиста", "success");
      fetchPlaylistTracks(selectedPlaylist.playlist_id);
    } catch { showSnackbar("Ошибка удаления из плейлиста", "error"); }
  };

  const showSnackbar = (message, severity) => {
    setSnackbar({ open: true, message, severity });
  };

  const handleTabChange = (e, val) => {
    setTab(val);
    if (val === 2) fetchStats();
    if (val === 3) fetchPlaylists();
  };

  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <AppBar position="static" sx={{ background: "linear-gradient(90deg, #1a0533 0%, #0a0a0a 100%)", boxShadow: "0 2px 20px rgba(124,77,255,0.3)" }}>
        <Toolbar>
          <LibraryMusicIcon sx={{ mr: 1, color: "#7c4dff" }} />
          <Typography variant="h6" sx={{ flexGrow: 1, fontWeight: 700, letterSpacing: 1 }}>
            Music Library
          </Typography>
          <Chip label={`${tracks.length} треков`} color="primary" size="small" />
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg" sx={{ mt: 3 }}>
        <Tabs value={tab} onChange={handleTabChange} sx={{ mb: 3, "& .MuiTab-root": { fontWeight: 600 } }}>
          <Tab label="Коллекция" />
          <Tab label="Поиск" icon={<SearchIcon />} iconPosition="start" />
          <Tab label="Статистика" icon={<BarChartIcon />} iconPosition="start" />
          <Tab label="Плейлисты" icon={<QueueMusicIcon />} iconPosition="start" />
        </Tabs>

        {/* Коллекция */}
        {tab === 0 && (
          <Box>
            <Paper sx={{ p: 2, mb: 2, background: "#141414" }}>
              <Box sx={{ display: "flex", alignItems: "center", gap: 2, flexWrap: "wrap" }}>
                <FilterListIcon sx={{ color: "#7c4dff" }} />
                <Typography variant="body2" sx={{ color: "#aaa" }}>Фильтр по дате добавления:</Typography>
                <TextField size="small" label="От" type="date" value={filterFrom}
                  onChange={(e) => setFilterFrom(e.target.value)}
                  InputLabelProps={{ shrink: true }} sx={{ width: 160 }} />
                <TextField size="small" label="До" type="date" value={filterTo}
                  onChange={(e) => setFilterTo(e.target.value)}
                  InputLabelProps={{ shrink: true }} sx={{ width: 160 }} />
                <Button variant="contained" size="small" onClick={handleFilter}
                  sx={{ background: "linear-gradient(45deg, #7c4dff, #ff4081)" }}>
                  Применить
                </Button>
                {filterActive && (
                  <Button variant="outlined" size="small" onClick={handleClearFilter}
                    sx={{ borderColor: "#666", color: "#aaa" }}>
                    Сбросить
                  </Button>
                )}
                {filterActive && <Chip label={`${filteredTracks.length} треков`} size="small" color="secondary" />}
              </Box>
            </Paper>
            <Box sx={{ display: "flex", justifyContent: "flex-end", mb: 2 }}>
              <Button variant="contained" startIcon={<AddIcon />} onClick={() => setOpenAdd(true)}
                sx={{ background: "linear-gradient(45deg, #7c4dff, #ff4081)", fontWeight: 700 }}>
                Добавить трек
              </Button>
            </Box>
            <TrackTable tracks={filteredTracks} onDelete={handleDeleteTrack}
              onEdit={(t) => { setEditTrack(t); setOpenEdit(true); }}
              onAddToPlaylist={(t) => { setSelectedTrackForPlaylist(t); setOpenAddToPlaylist(true); }} />
          </Box>
        )}

        {/* Поиск */}
        {tab === 1 && (
          <Box>
            <Box sx={{ display: "flex", gap: 2, mb: 3 }}>
              <TextField fullWidth label="Исполнитель или название" variant="outlined"
                value={searchQuery} onChange={(e) => setSearchQuery(e.target.value)}
                onKeyPress={(e) => e.key === "Enter" && handleSearch()}
                sx={{ "& .MuiOutlinedInput-root": { "&.Mui-focused fieldset": { borderColor: "#7c4dff" } } }} />
              <Button variant="contained" onClick={handleSearch} startIcon={<SearchIcon />}
                sx={{ px: 3, background: "linear-gradient(45deg, #7c4dff, #ff4081)" }}>
                Найти
              </Button>
            </Box>
            {searchResults.length > 0 && (
              <TrackTable tracks={searchResults} onDelete={handleDeleteTrack}
                onEdit={(t) => { setEditTrack(t); setOpenEdit(true); }}
                onAddToPlaylist={(t) => { setSelectedTrackForPlaylist(t); setOpenAddToPlaylist(true); }} />
            )}
          </Box>
        )}

        {/* Статистика */}
        {tab === 2 && stats && (
          <Box>
            <Grid container spacing={3} sx={{ mb: 3 }}>
              <Grid size={{ xs: 12, md: 4 }}>
                <StatCard title="Всего треков" value={stats.total_tracks} color="#7c4dff" />
              </Grid>
              <Grid size={{ xs: 12, md: 4 }}>
                <StatCard title="Исполнителей" value={Object.keys(stats.artist_counts || {}).length} color="#ff4081" />
              </Grid>
              <Grid size={{ xs: 12, md: 4 }}>
                <StatCard title="Топ исполнитель" value={stats.popular_artist || "—"} color="#00bcd4" />
              </Grid>
            </Grid>
            <Paper sx={{ p: 3, background: "#141414" }}>
              <Typography variant="h6" sx={{ mb: 2, fontWeight: 700 }}>Треки по исполнителям</Typography>
              {Object.entries(stats.artist_counts || {}).map(([artist, count]) => (
                <Box key={artist} sx={{ mb: 1.5 }}>
                  <Box sx={{ display: "flex", justifyContent: "space-between", mb: 0.5 }}>
                    <Typography variant="body2">{artist}</Typography>
                    <Typography variant="body2" color="primary">{count}</Typography>
                  </Box>
                  <Box sx={{ height: 6, borderRadius: 3, background: "#1e1e1e" }}>
                    <Box sx={{
                      height: "100%", borderRadius: 3,
                      width: `${(count / stats.total_tracks) * 100}%`,
                      background: "linear-gradient(90deg, #7c4dff, #ff4081)",
                      transition: "width 0.5s ease"
                    }} />
                  </Box>
                </Box>
              ))}
            </Paper>
          </Box>
        )}

        {/* Плейлисты */}
        {tab === 3 && (
          <Grid container spacing={3}>
            <Grid size={{ xs: 12, md: 4 }}>
              <Paper sx={{ background: "#141414", borderRadius: 2, p: 2 }}>
                <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2 }}>
                  <Typography variant="h6" sx={{ fontWeight: 700 }}>Плейлисты</Typography>
                  <IconButton onClick={() => setOpenAddPlaylist(true)} sx={{ color: "#7c4dff" }}>
                    <AddIcon />
                  </IconButton>
                </Box>
                <List disablePadding>
                  {playlists.map((pl, idx) => (
                    <Box key={pl.playlist_id}>
                      <ListItem button selected={selectedPlaylist?.playlist_id === pl.playlist_id}
                        onClick={() => { setSelectedPlaylist(pl); fetchPlaylistTracks(pl.playlist_id); }}
                        sx={{ borderRadius: 1, "&.Mui-selected": { background: "#7c4dff22" }, "&:hover": { background: "#1e1e1e" } }}>
                        <QueueMusicIcon sx={{ mr: 1.5, color: "#7c4dff", fontSize: 20 }} />
                        <ListItemText primary={pl.name} />
                        <ListItemSecondaryAction>
                          <IconButton size="small" onClick={(e) => { e.stopPropagation(); setSelectedPlaylist(pl); setRenameValue(pl.name); setOpenRenamePlaylist(true); }} sx={{ color: "#7c4dff" }}>
                            <EditIcon fontSize="small" />
                          </IconButton>
                          <IconButton size="small" onClick={(e) => { e.stopPropagation(); handleDeletePlaylist(pl.playlist_id); }} sx={{ color: "#ff4081" }}>
                            <DeleteIcon fontSize="small" />
                          </IconButton>
                        </ListItemSecondaryAction>
                      </ListItem>
                      {idx < playlists.length - 1 && <Divider sx={{ borderColor: "#1e1e1e" }} />}
                    </Box>
                  ))}
                  {playlists.length === 0 && (
                    <Typography variant="body2" sx={{ color: "#666", textAlign: "center", py: 3 }}>
                      Нет плейлистов
                    </Typography>
                  )}
                </List>
              </Paper>
            </Grid>
            <Grid item xs={12} md={8}>
              <Paper sx={{ background: "#141414", borderRadius: 2, p: 2 }}>
                {selectedPlaylist ? (
                  <>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 2 }}>
                      {selectedPlaylist.name}
                      <Chip label={`${playlistTracks.length} треков`} size="small" color="primary" sx={{ ml: 1 }} />
                    </Typography>
                    <Table size="small">
                      <TableHead>
                        <TableRow sx={{ "& th": { color: "#7c4dff", fontWeight: 700 } }}>
                          <TableCell>Исполнитель</TableCell>
                          <TableCell>Название</TableCell>
                          <TableCell>Ссылка</TableCell>
                          <TableCell align="right">Удалить</TableCell>
                        </TableRow>
                      </TableHead>
                      <TableBody>
                        {playlistTracks.map((track) => (
                          <TableRow key={track.track_id} sx={{ "&:hover": { background: "#1e1e1e" } }}>
                            <TableCell sx={{ fontWeight: 600 }}>{track.artist}</TableCell>
                            <TableCell>{track.track}</TableCell>
                            <TableCell>
                              {track.url ? (
                                <a href={track.url} target="_blank" rel="noreferrer" style={{ color: "#7c4dff", textDecoration: "none" }}>
                                  Открыть
                                </a>
                              ) : "—"}
                            </TableCell>
                            <TableCell align="right">
                              <IconButton size="small" onClick={() => handleRemoveFromPlaylist(track.track_id)} sx={{ color: "#ff4081" }}>
                                <DeleteIcon fontSize="small" />
                              </IconButton>
                            </TableCell>
                          </TableRow>
                        ))}
                        {playlistTracks.length === 0 && (
                          <TableRow>
                            <TableCell colSpan={4} align="center" sx={{ color: "#666", py: 3 }}>
                              Плейлист пуст — добавьте треки из коллекции
                            </TableCell>
                          </TableRow>
                        )}
                      </TableBody>
                    </Table>
                  </>
                ) : (
                  <Typography variant="body2" sx={{ color: "#666", textAlign: "center", py: 5 }}>
                    Выберите плейлист слева
                  </Typography>
                )}
              </Paper>
            </Grid>
          </Grid>
        )}
      </Container>

      {/* Диалог добавления трека */}
      <Dialog open={openAdd} onClose={() => setOpenAdd(false)} PaperProps={{ sx: { background: "#141414", minWidth: 400 } }}>
        <DialogTitle sx={{ fontWeight: 700 }}>Добавить трек</DialogTitle>
        <DialogContent sx={{ display: "flex", flexDirection: "column", gap: 2, pt: 2 }}>
          <TextField label="Исполнитель" value={newTrack.artist} onChange={(e) => setNewTrack({ ...newTrack, artist: e.target.value })} fullWidth />
          <TextField label="Название трека" value={newTrack.track} onChange={(e) => setNewTrack({ ...newTrack, track: e.target.value })} fullWidth />
          <TextField label="Ссылка на клип" value={newTrack.url} onChange={(e) => setNewTrack({ ...newTrack, url: e.target.value })} fullWidth />
        </DialogContent>
        <DialogActions sx={{ p: 2 }}>
          <Button onClick={() => setOpenAdd(false)}>Отмена</Button>
          <Button variant="contained" onClick={handleAddTrack} sx={{ background: "linear-gradient(45deg, #7c4dff, #ff4081)" }}>Добавить</Button>
        </DialogActions>
      </Dialog>

      {/* Диалог редактирования трека */}
      <Dialog open={openEdit} onClose={() => setOpenEdit(false)} PaperProps={{ sx: { background: "#141414", minWidth: 400 } }}>
        <DialogTitle sx={{ fontWeight: 700 }}>Редактировать трек</DialogTitle>
        <DialogContent sx={{ display: "flex", flexDirection: "column", gap: 2, pt: 2 }}>
          <TextField label="Исполнитель" value={editTrack.artist} onChange={(e) => setEditTrack({ ...editTrack, artist: e.target.value })} fullWidth />
          <TextField label="Название трека" value={editTrack.track} onChange={(e) => setEditTrack({ ...editTrack, track: e.target.value })} fullWidth />
          <TextField label="Ссылка на клип" value={editTrack.url} onChange={(e) => setEditTrack({ ...editTrack, url: e.target.value })} fullWidth />
        </DialogContent>
        <DialogActions sx={{ p: 2 }}>
          <Button onClick={() => setOpenEdit(false)}>Отмена</Button>
          <Button variant="contained" onClick={handleEditTrack} sx={{ background: "linear-gradient(45deg, #7c4dff, #ff4081)" }}>Сохранить</Button>
        </DialogActions>
      </Dialog>

      {/* Диалог создания плейлиста */}
      <Dialog open={openAddPlaylist} onClose={() => setOpenAddPlaylist(false)} PaperProps={{ sx: { background: "#141414", minWidth: 350 } }}>
        <DialogTitle sx={{ fontWeight: 700 }}>Новый плейлист</DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          <TextField label="Название плейлиста" value={newPlaylistName}
            onChange={(e) => setNewPlaylistName(e.target.value)} fullWidth
            onKeyPress={(e) => e.key === "Enter" && handleCreatePlaylist()} />
        </DialogContent>
        <DialogActions sx={{ p: 2 }}>
          <Button onClick={() => setOpenAddPlaylist(false)}>Отмена</Button>
          <Button variant="contained" onClick={handleCreatePlaylist} sx={{ background: "linear-gradient(45deg, #7c4dff, #ff4081)" }}>Создать</Button>
        </DialogActions>
      </Dialog>

      {/* Диалог переименования плейлиста */}
      <Dialog open={openRenamePlaylist} onClose={() => setOpenRenamePlaylist(false)} PaperProps={{ sx: { background: "#141414", minWidth: 350 } }}>
        <DialogTitle sx={{ fontWeight: 700 }}>Переименовать плейлист</DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          <TextField label="Новое название" value={renameValue}
            onChange={(e) => setRenameValue(e.target.value)} fullWidth
            onKeyPress={(e) => e.key === "Enter" && handleRenamePlaylist()} />
        </DialogContent>
        <DialogActions sx={{ p: 2 }}>
          <Button onClick={() => setOpenRenamePlaylist(false)}>Отмена</Button>
          <Button variant="contained" onClick={handleRenamePlaylist} sx={{ background: "linear-gradient(45deg, #7c4dff, #ff4081)" }}>Сохранить</Button>
        </DialogActions>
      </Dialog>

      {/* Диалог добавления в плейлист */}
      <Dialog open={openAddToPlaylist} onClose={() => setOpenAddToPlaylist(false)} PaperProps={{ sx: { background: "#141414", minWidth: 350 } }}>
        <DialogTitle sx={{ fontWeight: 700 }}>
          Добавить в плейлист
          {selectedTrackForPlaylist && (
            <Typography variant="body2" sx={{ color: "#aaa", mt: 0.5 }}>
              {selectedTrackForPlaylist.artist} — {selectedTrackForPlaylist.track}
            </Typography>
          )}
        </DialogTitle>
        <DialogContent>
          <List disablePadding>
            {playlists.map((pl) => (
              <ListItem button key={pl.playlist_id} onClick={() => handleAddTrackToPlaylist(pl.playlist_id)}
                sx={{ borderRadius: 1, "&:hover": { background: "#7c4dff22" } }}>
                <QueueMusicIcon sx={{ mr: 1.5, color: "#7c4dff", fontSize: 20 }} />
                <ListItemText primary={pl.name} />
              </ListItem>
            ))}
            {playlists.length === 0 && (
              <Typography variant="body2" sx={{ color: "#666", textAlign: "center", py: 2 }}>
                Сначала создайте плейлист
              </Typography>
            )}
          </List>
        </DialogContent>
        <DialogActions sx={{ p: 2 }}>
          <Button onClick={() => setOpenAddToPlaylist(false)}>Отмена</Button>
        </DialogActions>
      </Dialog>

      <Snackbar open={snackbar.open} autoHideDuration={3000} onClose={() => setSnackbar({ ...snackbar, open: false })}>
        <Alert severity={snackbar.severity}>{snackbar.message}</Alert>
      </Snackbar>
    </ThemeProvider>
  );
}

function TrackTable({ tracks, onDelete, onEdit, onAddToPlaylist }) {
  return (
    <TableContainer component={Paper} sx={{ background: "#141414", borderRadius: 2 }}>
      <Table>
        <TableHead>
          <TableRow sx={{ "& th": { fontWeight: 700, color: "#7c4dff", borderBottom: "2px solid #7c4dff22" } }}>
            <TableCell>ID</TableCell>
            <TableCell>Исполнитель</TableCell>
            <TableCell>Название</TableCell>
            <TableCell>Добавлен</TableCell>
            <TableCell>Ссылка</TableCell>
            <TableCell>YouTube</TableCell>
            <TableCell align="right">Действия</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {tracks.map((track) => (
            <TableRow key={track.track_id} sx={{ "&:hover": { background: "#1e1e1e" }, transition: "background 0.2s" }}>
              <TableCell sx={{ color: "#666" }}>{track.track_id}</TableCell>
              <TableCell sx={{ fontWeight: 600 }}>{track.artist}</TableCell>
              <TableCell>{track.track}</TableCell>
              <TableCell sx={{ color: "#888", fontSize: 12 }}>
                {track.created_at ? new Date(track.created_at).toLocaleDateString("ru-RU") : "—"}
              </TableCell>
              <TableCell>
                {track.url ? (
                  <a href={track.url} target="_blank" rel="noreferrer" style={{ color: "#7c4dff", textDecoration: "none" }}>
                    Открыть
                  </a>
                ) : "—"}
              </TableCell>
              <TableCell>
                <IconButton size="small" onClick={() => window.open(
                  `https://www.youtube.com/results?search_query=${encodeURIComponent(track.artist + " " + track.track)}`, "_blank"
                )}>
                  <svg viewBox="0 0 24 24" width="20" height="20" fill="#ff0000">
                    <path d="M23.5 6.2s-.3-2-1.2-2.8c-1.1-1.2-2.4-1.2-3-1.3C16.8 2 12 2 12 2s-4.8 0-7.3.1c-.6.1-1.9.1-3 1.3C.8 4.2.5 6.2.5 6.2S.2 8.5.2 10.8v2.1c0 2.3.3 4.6.3 4.6s.3 2 1.2 2.8c1.1 1.2 2.6 1.1 3.3 1.2C7.2 21.7 12 21.8 12 21.8s4.8 0 7.3-.2c.6-.1 1.9-.1 3-1.3.9-.8 1.2-2.8 1.2-2.8s.3-2.3.3-4.6v-2.1c0-2.3-.3-4.6-.3-4.6zM9.7 15.5V8.4l8.1 3.6-8.1 3.5z" />
                  </svg>
                </IconButton>
              </TableCell>
              <TableCell align="right">
                <IconButton onClick={() => onAddToPlaylist(track)} size="small" sx={{ color: "#00bcd4", mr: 0.5 }}>
                  <PlaylistAddIcon fontSize="small" />
                </IconButton>
                <IconButton onClick={() => onEdit(track)} size="small" sx={{ color: "#7c4dff", mr: 0.5 }}>
                  <EditIcon fontSize="small" />
                </IconButton>
                <IconButton onClick={() => onDelete(track.track_id)} size="small" sx={{ color: "#ff4081" }}>
                  <DeleteIcon fontSize="small" />
                </IconButton>
              </TableCell>
            </TableRow>
          ))}
          {tracks.length === 0 && (
            <TableRow>
              <TableCell colSpan={7} align="center" sx={{ color: "#666", py: 4 }}>
                Треки не найдены
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </TableContainer>
  );
}

function StatCard({ title, value, color }) {
  return (
    <Card sx={{ background: "#141414", borderRadius: 2, border: `1px solid ${color}33` }}>
      <CardContent>
        <Typography variant="body2" sx={{ color: "#666", mb: 1 }}>{title}</Typography>
        <Typography variant="h4" sx={{ fontWeight: 700, color }}>{value}</Typography>
      </CardContent>
    </Card>
  );
}