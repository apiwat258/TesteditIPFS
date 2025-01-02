const express = require('express');
const { Pool } = require('pg');
const cors = require('cors');
const path = require('path'); // Import path module for serving static files
const axios = require('axios'); // Import axios for making HTTP requests

const app = express();
const port = 3001;

// Middleware
app.use(cors());
app.use(express.json());

// Serve static files from the current directory
app.use(express.static(__dirname));

// PostgreSQL connection
const pool = new Pool({
    user: 'postgres',
    host: 'localhost',
    database: 'postgres',
    password: 'password',
    port: 5432,
});

// Default route serving displaymain.html
app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, 'displaymain.html'));
});

// Route for edit-page.html
app.get('/edit-page', (req, res) => {
    res.sendFile(path.join(__dirname, 'edit-page.html'));
});

// Route for playIPFS.html
app.get('/playIPFS.html', (req, res) => {
    res.sendFile(path.join(__dirname, 'playIPFS.html'));
});

app.get('/playPostgres.html', (req, res) => {
    res.sendFile(path.join(__dirname, 'playPostgres.html'));
});

// GET endpoint to fetch data from only Postgresql from ID
app.get('/raw-milk-data/:id', async (req, res) => {
    const { id } = req.params; // Get the id from the URL
    try {
        const result = await pool.query('SELECT * FROM public."BatchRawMilk" WHERE id = $1', [id]);
        if (result.rows.length === 0) {
            return res.status(404).json({ error: `No data found for id ${id}.` });
        }
        res.status(200).json(result.rows[0]); // Send the matching row as the response
    } catch (error) {
        console.error('Error fetching data from PostgreSQL: ', error);
        res.status(500).json({ error: 'Internal Server Error' });
    }
});



// Route to fetch milk batch data
app.get('/milk-batches', async (req, res) => {
    try {
        const query = `
            SELECT 
                id,
                person_in_charge, 
                milk_tank_number, 
                updated_on, 
                ipfs_cid 
            FROM public."BatchRawMilk";
        `;
        const result = await pool.query(query);
        res.json(result.rows);
    } catch (err) {
        console.error(err);
        res.status(500).send('Server error');
    }
});

// Start the server
app.listen(port, () => {
    console.log(`Server running on http://localhost:${port}`);
});
