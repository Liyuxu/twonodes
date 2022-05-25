import argparse

# GLOBAL PARAMETERS

# SERVER_ADDR= 'localhost'   # When running in a real distributed setting, change to the server's IP address
SERVER_ADDR = '127.0.0.1'  # When running in a real distributed setting, change to the server's IP address
SERVER_PORT = 50000
ClientID = 0


def read_options():
    parser = argparse.ArgumentParser()

    parser.add_argument('--dataset',
                        help='name of dataset;',
                        type=str,
                        default='mnist')
    parser.add_argument('--model',
                        help='name of model;',
                        type=str,
                        default='logistic')
    parser.add_argument('--gpu',
                        action='store_true',
                        default=False,
                        help='use gpu (default: False)')
    parser.add_argument('--num_round',
                        help='number of rounds to simulate;',
                        type=int,
                        default=10000)
    parser.add_argument('--clients_per_round',
                        help='number of clients trained per round;',
                        type=int,
                        default=2)
    parser.add_argument('--batch_size',
                        help='batch size when clients train on data;',
                        type=int,
                        default=5000)
    parser.add_argument('--lr',
                        help='learning rate for inner solver;',
                        type=float,
                        default=1)
    parser.add_argument('--out_dim',
                        help='output dimension',
                        type=int,
                        default=10)

    parsed = parser.parse_args()
    options = parsed.__dict__
    options['gpu'] = options['gpu'] and torch.cuda.is_available()

    return options
